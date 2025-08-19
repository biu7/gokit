package amqp

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/biu7/gokit/log"

	"github.com/wagslane/go-rabbitmq"
)

type Publisher struct {
	shutdown bool

	conn       *rabbitmq.Conn
	publishers *rabbitmq.Publisher
	mutex      sync.Mutex
	log        rabbitmq.Logger
}

func NewPublisher(conf *Config) (*Publisher, error) {
	logger := NewLogger(log.Default)
	conn, err := rabbitmq.NewConn(
		fmt.Sprintf("amqp://%s:%s@%s/%s", conf.Username, conf.Password, conf.Endpoint, conf.Vhost),
		rabbitmq.WithConnectionOptionsLogging,
		rabbitmq.WithConnectionOptionsLogger(logger),
	)
	if err != nil {
		return nil, fmt.Errorf("new rabbitmq publisher failed: %w", err)
	}

	pub := &Publisher{
		conn:  conn,
		mutex: sync.Mutex{},
		log:   logger,
	}
	return pub, nil
}

func (r *Publisher) getExchangePublisher() (*rabbitmq.Publisher, error) {
	if r.shutdown {
		return nil, errors.New("producer is shutdown")
	}
	if r.publishers != nil {
		return r.publishers, nil
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.publishers != nil {
		return r.publishers, nil
	}

	var opts = []func(*rabbitmq.PublisherOptions){
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsLogger(r.log),
	}
	publisher, err := rabbitmq.NewPublisher(
		r.conn,
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("new publisher failed: %w", err)
	}

	publisher.NotifyPublish(r.notifyPublish)
	publisher.NotifyReturn(r.notifyReturn)
	r.publishers = publisher
	return publisher, nil
}

func (r *Publisher) Publish(data []byte, exchange string, keys []string, delay time.Duration) error {
	publisher, err := r.getExchangePublisher()
	if err != nil {
		return err
	}
	var opts = []func(*rabbitmq.PublishOptions){
		rabbitmq.WithPublishOptionsTimestamp(time.Now()),
		rabbitmq.WithPublishOptionsExchange(exchange),
	}
	if delay > 0 {
		opts = append(opts, rabbitmq.WithPublishOptionsHeaders(rabbitmq.Table{
			"delay": delay.Milliseconds(),
		}))
	}

	return publisher.Publish(data, keys, opts...)
}

func (r *Publisher) Close() {
	r.shutdown = true
	if r.publishers != nil {
		r.publishers.Close()
	}
	// r.conn.Close()
}

func (r *Publisher) notifyPublish(p rabbitmq.Confirmation) {
	// log.Info("[AMQP] NotifyPublish", "de", p.DeliveryTag, "ack", p.Ack, "rc", p.ReconnectionCount)
}

func (r *Publisher) notifyReturn(p rabbitmq.Return) {
	log.Warn("[AMQP] NotifyReturn", "code", p.ReplyCode, "reason", p.ReplyText, "exchange", p.Exchange, "routing_key", p.RoutingKey, "body", string(p.Body))
}
