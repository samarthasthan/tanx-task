package rabbitmq

import (
	ampq "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	*ampq.Connection
}

func NewRabbitMQ(a string) (*RabbitMQ, error) {
	conn, err := ampq.Dial(a)
	if err != nil {
		return nil, err
	}
	return &RabbitMQ{conn}, nil
}
