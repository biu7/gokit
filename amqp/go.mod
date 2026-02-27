module github.com/biu7/gokit/amqp

go 1.23.0

require (
	github.com/biu7/gokit v0.0.0
	github.com/wagslane/go-rabbitmq v0.13.0
)

require (
	github.com/rabbitmq/amqp091-go v1.7.0 // indirect
	go.opentelemetry.io/otel v1.24.0 // indirect
	go.opentelemetry.io/otel/trace v1.24.0 // indirect
)

replace github.com/biu7/gokit => ../
