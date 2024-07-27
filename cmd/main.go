package main

import (
	"fmt"
	"log"
	"time"

	"github.com/samarthasthan/tanx-task/internal/database"
	"github.com/samarthasthan/tanx-task/internal/rabbitmq"
	"github.com/samarthasthan/tanx-task/pkg/env"
)

// Define the environment variables
var (
	REST_API_PORT         string
	MYSQL_PORT            string
	MYSQL_ROOT_PASSWORD   string
	MYSQL_HOST            string
	REDIS_PORT            string
	REDIS_HOST            string
	RABBITMQ_DEFAULT_USER string
	RABBITMQ_DEFAULT_PASS string
	RABBITMQ_DEFAULT_PORT string
	RABBITMQ_DEFAULT_HOST string
	SMTP_SERVER           string
	SMTP_PORT             string
	SMTP_LOGIN            string
	SMTP_PASSWORD         string
)

// Initialize the environment variables
func init() {
	REST_API_PORT = env.GetEnv("REST_API_PORT", "3000")
	MYSQL_PORT = env.GetEnv("MYSQL_PORT", "3306")
	MYSQL_ROOT_PASSWORD = env.GetEnv("MYSQL_ROOT_PASSWORD", "password")
	MYSQL_HOST = env.GetEnv("MYSQL_HOST", "localhost")
	REDIS_PORT = env.GetEnv("REDIS_PORT", "6379")
	REDIS_HOST = env.GetEnv("REDIS_HOST", "localhost")
	RABBITMQ_DEFAULT_USER = env.GetEnv("RABBITMQ_DEFAULT_USER", "root")
	RABBITMQ_DEFAULT_PASS = env.GetEnv("RABBITMQ_DEFAULT_PASS", "password")
	RABBITMQ_DEFAULT_PORT = env.GetEnv("RABBITMQ_DEFAULT_PORT", "5672")
	RABBITMQ_DEFAULT_HOST = env.GetEnv("RABBITMQ_DEFAULT_HOST", "localhost")
	SMTP_SERVER = env.GetEnv("SMTP_SERVER", "smtp-relay.sendinblue.com")
	SMTP_PORT = env.GetEnv("SMTP_PORT", "587")
	SMTP_LOGIN = env.GetEnv("SMTP_LOGIN", "use your own sender")
	SMTP_PASSWORD = env.GetEnv("SMTP_PASSWORD", "use your own key")
}

func main() {
	// Create a new instance of MySQL database
	mysql := database.NewMySQL()
	err := mysql.Connect(fmt.Sprintf("root:%s@tcp(%s:%s)/tanx?parseTime=true", MYSQL_ROOT_PASSWORD, MYSQL_HOST, MYSQL_PORT))
	if err != nil {
		log.Println(err.Error())
	}

	// We can use same instance of RabbitMQ for both publisher and consumer
	// But for the sake of demonstration, we will create two instances
	// As both publisher and consumer will be running in the different containers
	// Create a new RabbitMQ instance for the publisher
	publisher, err := rabbitmq.NewRabbitMQ(fmt.Sprintf("amqp://%s:%s@%s:%s/", RABBITMQ_DEFAULT_USER, RABBITMQ_DEFAULT_PASS, RABBITMQ_DEFAULT_HOST, RABBITMQ_DEFAULT_PORT))
	if err != nil {
		failOnError(err, "Failed to connect to RabbitMQ as publisher")
	}
	defer publisher.Close()
	go func() {
		for {
			// Publish a message to the RabbitMQ
			err := publisher.Publish("tanx", "tanx", []byte("Hello World"))
			if err != nil {
				log.Println(err.Error())
			}
			time.Sleep(time.Second)
		}
	}()

	// Create a new RabbitMQ instance for the consumer
	consumer, err := rabbitmq.NewRabbitMQ(fmt.Sprintf("amqp://%s:%s@%s:%s/", RABBITMQ_DEFAULT_USER, RABBITMQ_DEFAULT_PASS, RABBITMQ_DEFAULT_HOST, RABBITMQ_DEFAULT_PORT))
	if err != nil {
		failOnError(err, "Failed to connect to RabbitMQ as consumer")
	}
	defer consumer.Close()

	for {
		// Consume a message from the RabbitMQ
		msgs, err := consumer.Consume("tanx", "tanx")
		if err != nil {
			log.Println(err.Error())
		}
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}

	// // Register controllers
	// c := controller.NewController(publisher, mysql)

	// // Register Echo handler
	// h := api.NewHandler(c)
	// h.Handle()
	// h.Logger.Fatal(h.Start(":" + REST_API_PORT))
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
