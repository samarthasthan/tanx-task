package main

import (
	"fmt"

	"github.com/samarthasthan/tanx-task/internal/database"
	"github.com/samarthasthan/tanx-task/internal/rabbitmq"
	"github.com/samarthasthan/tanx-task/pkg/env"
	rabbitmq_utils "github.com/samarthasthan/tanx-task/pkg/rabbitmq"
)

// Define the environment variables
var (
	REDIS_PORT            string
	REDIS_HOST            string
	RABBITMQ_DEFAULT_USER string
	RABBITMQ_DEFAULT_PASS string
	RABBITMQ_DEFAULT_PORT string
	RABBITMQ_DEFAULT_HOST string
)

// Initialize the environment variables
func init() {
	REDIS_PORT = env.GetEnv("REDIS_PORT", "6379")
	REDIS_HOST = env.GetEnv("REDIS_HOST", "localhost")
	RABBITMQ_DEFAULT_USER = env.GetEnv("RABBITMQ_DEFAULT_USER", "root")
	RABBITMQ_DEFAULT_PASS = env.GetEnv("RABBITMQ_DEFAULT_PASS", "password")
	RABBITMQ_DEFAULT_PORT = env.GetEnv("RABBITMQ_DEFAULT_PORT", "5672")
	RABBITMQ_DEFAULT_HOST = env.GetEnv("RABBITMQ_DEFAULT_HOST", "localhost")
}

func main() {
	// Create a new instance of Redis database
	redis := database.NewRedis()
	err := redis.Connect(fmt.Sprintf("%s:%s", REDIS_HOST, REDIS_PORT))
	if err != nil {
		panic(err.Error())
	}
	defer redis.Close()

	// We can use same instance of RabbitMQ for both publisher and consumer
	// As both publisher and consumer will be running in the different containers
	// Create a new RabbitMQ instance for the publisher
	publisher, err := rabbitmq.NewRabbitMQ(fmt.Sprintf("amqp://%s:%s@%s:%s/", RABBITMQ_DEFAULT_USER, RABBITMQ_DEFAULT_PASS, RABBITMQ_DEFAULT_HOST, RABBITMQ_DEFAULT_PORT))
	if err != nil {
		rabbitmq_utils.FailOnError(err, "Failed to connect to RabbitMQ as publisher")
	}

	defer publisher.Close()
}
