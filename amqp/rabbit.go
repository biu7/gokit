package amqp

type Rabbit struct {
	publisher *Publisher
	consumer  *Consumer
}

func NewRabbit(conf *Config) (*Rabbit, func(), error) {
	publisher, err := NewPublisher(conf)
	if err != nil {
		return nil, nil, err
	}
	consumer, err := NewConsumer(conf)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		publisher.Close()
		consumer.Close()
	}
	return &Rabbit{
		publisher: publisher,
		consumer:  consumer,
	}, cleanup, nil
}

func (r *Rabbit) Consumer() *Consumer {
	return r.consumer
}

func (r *Rabbit) Publisher() *Publisher {
	return r.publisher
}

type Config struct {
	Username string
	Password string
	Endpoint string
	Vhost    string
}
