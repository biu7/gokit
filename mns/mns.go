package mns

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	mns "github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/biu7/gokit/log"
	"github.com/biu7/gokit/safe"
	"strings"
)

type QueueMessage[T any] struct {
	Body T `json:"body"`
}

type ReceiveMessage struct {
	ID   string `json:"id"`
	Body string `json:"body"`
}

type Queue struct {
	name     string
	shutdown bool
	queue    mns.AliMNSQueue
}

func NewQueue(queueName, endpoint, accessKeyID, accessKeySecret string) *Queue {
	client := mns.NewAliMNSClientWithConfig(mns.AliMNSClientConfig{
		EndPoint:        endpoint,
		AccessKeyId:     accessKeyID,
		AccessKeySecret: accessKeySecret,
		TimeoutSecond:   mns.DefaultTimeout,
		MaxConnsPerHost: mns.DefaultMaxConnsPerHost,
	})
	queue := mns.NewMNSQueue(queueName, client)
	return &Queue{queue: queue, name: queueName}
}

func (q *Queue) Name() string {
	return q.name
}

func (q *Queue) Send(msg *QueueMessage[any]) (string, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("could not marshal message: %w", err)
	}
	result, err := q.queue.SendMessage(mns.MessageSendRequest{
		MessageBody: base64.StdEncoding.EncodeToString(data),
		Priority:    10,
	})
	if err != nil {
		return "", fmt.Errorf("could not send message: %w", err)
	}
	return result.MessageId, nil
}

func (q *Queue) Consume(ctx context.Context, f func(ReceiveMessage) error) {
	respChan := make(chan mns.MessageReceiveResponse)
	errChan := make(chan error)
	safe.Go(context.Background(), func(ctx context.Context) {
		q.processReceive(ctx, respChan, errChan, f)
	}, log.Default)
	log.Default.Info("[MNS] consume mns queue start", "queue", q.name)
	// 接收消息
	q.receive(ctx, respChan, errChan)
}

func (q *Queue) processReceive(ctx context.Context, respChan chan mns.MessageReceiveResponse, errChan chan error, f func(ReceiveMessage) error) {
	for {
		select {
		case resp := <-respChan:
			err := f(ReceiveMessage{
				ID:   resp.MessageId,
				Body: resp.MessageBody,
			})
			if err != nil {
				log.Error("[MNS] consume mns msg error", "error", err, "msgId", resp.MessageId, "queue", q.name)
				continue
			}
			if err = q.queue.DeleteMessage(resp.ReceiptHandle); err != nil {
				log.Error("delete mns message error", "error", err, "msgId", resp.MessageId, "queue", q.name)
				continue
			}
		case err := <-errChan:
			if strings.Contains(err.Error(), "MessageNotExist") {
				continue
			}
			log.Error("receive mns message error", "error", err, "queue", q.name)
		case <-ctx.Done():
			return
		}
	}
}

func (q *Queue) receive(ctx context.Context, respChan chan mns.MessageReceiveResponse, errChan chan error) {
	for {
		q.queue.ReceiveMessage(respChan, errChan, 30)
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}
