package main

import (
	"fmt"
	"log"

	"github.com/samarthasthan/tanx-task/api"
	"github.com/samarthasthan/tanx-task/internal/controller"
	"github.com/samarthasthan/tanx-task/internal/database"
	"github.com/samarthasthan/tanx-task/internal/rabbitmq"
	"github.com/samarthasthan/tanx-task/pkg/env"
	rabbitmq_utils "github.com/samarthasthan/tanx-task/pkg/rabbitmq"
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
	JWT_SECRET            string
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
	SMTP_SERVER = env.GetEnv("SMTP_SERVER", "smtp-relay.brevo.com")
	SMTP_PORT = env.GetEnv("SMTP_PORT", "587")
	SMTP_LOGIN = env.GetEnv("SMTP_LOGIN", "75a33c001@smtp-brevo.com")
	SMTP_PASSWORD = env.GetEnv("SMTP_PASSWORD", "0c8shB9P4N3vXTyV")
	JWT_SECRET = env.GetEnv("JWT_SECRET", "secret")
}

func main() {
	// Create a new instance of MySQL database
	mysql := database.NewMySQL()
	err := mysql.Connect(fmt.Sprintf("root:%s@tcp(%s:%s)/tanx?parseTime=true", MYSQL_ROOT_PASSWORD, MYSQL_HOST, MYSQL_PORT))
	if err != nil {
		log.Println(err.Error())
	}

	// We can use same instance of RabbitMQ for both publisher and consumer
	// As both publisher and consumer will be running in the different containers
	// Create a new RabbitMQ instance for the publisher
	publisher, err := rabbitmq.NewRabbitMQ(fmt.Sprintf("amqp://%s:%s@%s:%s/", RABBITMQ_DEFAULT_USER, RABBITMQ_DEFAULT_PASS, RABBITMQ_DEFAULT_HOST, RABBITMQ_DEFAULT_PORT))
	if err != nil {
		rabbitmq_utils.FailOnError(err, "Failed to connect to RabbitMQ as publisher")
	}
	defer publisher.Close()

	// Register controllers
	c := controller.NewController(publisher, mysql, JWT_SECRET)

	// Register Echo handler
	h := api.NewHandler(c)
	h.Handle()
	h.Logger.Fatal(h.Start(":" + REST_API_PORT))
}
