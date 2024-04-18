package amqp

import (
	"fmt"
	"os"

	"github.com/biu7/gokit-qi/log"

	"github.com/wagslane/go-rabbitmq"
)

type Consumer struct {
	conn      *rabbitmq.Conn
	consumers []*rabbitmq.Consumer
	log       rabbitmq.Logger
}

func NewConsumer(conf *Config) (*Consumer, error) {
	logger := NewLogger(log.Default)

	conn, err := rabbitmq.NewConn(
		fmt.Sprintf("amqp://%s:%s@%s/%s", conf.Username, conf.Password, conf.Endpoint, conf.Vhost),
		rabbitmq.WithConnectionOptionsLogging,
		rabbitmq.WithConnectionOptionsLogger(logger),
	)
	if err != nil {
		return nil, fmt.Errorf("new rabbitmq consumer failed: %w", err)
	}

	cs := &Consumer{
		conn: conn,
		log:  logger,
	}

	return cs, nil
}

func (c *Consumer) Consume(exchange, queue string, keys []string, handler rabbitmq.Handler) error {
	hostname, _ := os.Hostname()
	var opts = []func(*rabbitmq.ConsumerOptions){
		rabbitmq.WithConsumerOptionsConcurrency(10),
		rabbitmq.WithConsumerOptionsConsumerName(hostname),
		rabbitmq.WithConsumerOptionsQueueDurable,

		rabbitmq.WithConsumerOptionsExchangeName(exchange),
		rabbitmq.WithConsumerOptionsLogging,
		rabbitmq.WithConsumerOptionsLogger(c.log),
	}
	for _, key := range keys {
		opts = append(opts, rabbitmq.WithConsumerOptionsRoutingKey(key))
	}

	consumer, err := rabbitmq.NewConsumer(
		c.conn,
		queue,
		opts...,
	)
	if err != nil {
		return fmt.Errorf("new consumer failed: %w", err)
	}
	c.consumers = append(c.consumers, consumer)
	return consumer.Run(handler)
}

func (c *Consumer) Close() {
	log.Info("[AMQP] consumer closing")
	for _, consumer := range c.consumers {
		consumer.Close()
	}
	// c.conn.Close()
	log.Info("[AMQP] consumer closed")
}
