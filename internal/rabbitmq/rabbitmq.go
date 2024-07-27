package rabbitmq

import (
	ampq "github.com/rabbitmq/amqp091-go"
)

// We can use queue interface to define the methods if we want to use different message brokers
// But for the sake of simplicity, we will use the RabbitMQ struct

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

func (r *RabbitMQ) Close() error {
	return r.Connection.Close()
}

func (r *RabbitMQ) Publish(exchange, key string, msg []byte) error {
	ch, err := r.Connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	if err := ch.Publish(exchange, key, false, false, ampq.Publishing{Body: msg}); err != nil {
		return err
	}
	return nil
}

func (r *RabbitMQ) Consume(queue, key string) (<-chan ampq.Delivery, error) {
	ch, err := r.Connection.Channel()
	if err != nil {
		return nil, err
	}
	if err := ch.ExchangeDeclare(queue, "direct", true, false, false, false, nil); err != nil {
		return nil, err
	}
	if _, err := ch.QueueDeclare(queue, true, false, false, false, nil); err != nil {
		return nil, err
	}
	if err := ch.QueueBind(queue, key, queue, false, nil); err != nil {
		return nil, err
	}
	msgs, err := ch.Consume(queue, "", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}
