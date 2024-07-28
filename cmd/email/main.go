package main

import (
	"fmt"
	"log"

	"github.com/samarthasthan/tanx-task/internal/mail"
	"github.com/samarthasthan/tanx-task/internal/rabbitmq"
	"github.com/samarthasthan/tanx-task/pkg/env"
	rabbitmq_utils "github.com/samarthasthan/tanx-task/pkg/rabbitmq"
)

var (
	RABBITMQ_DEFAULT_USER string
	RABBITMQ_DEFAULT_PASS string
	RABBITMQ_DEFAULT_PORT string
	RABBITMQ_DEFAULT_HOST string
	SMTP_SERVER           string
	SMTP_PORT             string
	SMTP_LOGIN            string
	SMTP_PASSWORD         string
)

func init() {
	RABBITMQ_DEFAULT_USER = env.GetEnv("RABBITMQ_DEFAULT_USER", "root")
	RABBITMQ_DEFAULT_PASS = env.GetEnv("RABBITMQ_DEFAULT_PASS", "password")
	RABBITMQ_DEFAULT_PORT = env.GetEnv("RABBITMQ_DEFAULT_PORT", "5672")
	RABBITMQ_DEFAULT_HOST = env.GetEnv("RABBITMQ_DEFAULT_HOST", "localhost")
	SMTP_SERVER = env.GetEnv("SMTP_SERVER", "smtp-relay.brevo.com")
	SMTP_PORT = env.GetEnv("SMTP_PORT", "587")
	SMTP_LOGIN = env.GetEnv("SMTP_LOGIN", "75a33c001@smtp-brevo.com")
	SMTP_PASSWORD = env.GetEnv("SMTP_PASSWORD", "0c8shB9P4N3vXTyV")
}

func main() {
	// Create a new RabbitMQ instance for the consumer
	consumer, err := rabbitmq.NewRabbitMQ(fmt.Sprintf("amqp://%s:%s@%s:%s/", RABBITMQ_DEFAULT_USER, RABBITMQ_DEFAULT_PASS, RABBITMQ_DEFAULT_HOST, RABBITMQ_DEFAULT_PORT))
	if err != nil {
		rabbitmq_utils.FailOnError(err, "Failed to connect to RabbitMQ as consumer")
	}
	defer consumer.Close()

	// Mail handler
	m := mail.NewMailHandler(consumer, SMTP_SERVER, SMTP_PORT, SMTP_LOGIN, SMTP_PASSWORD)
	log.Println("Consuming mails...")
	m.ConsumeMails()
}
