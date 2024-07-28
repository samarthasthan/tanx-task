package main

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	alert "github.com/samarthasthan/tanx-task/internal/alert"
	"github.com/samarthasthan/tanx-task/internal/database"
	"github.com/samarthasthan/tanx-task/internal/rabbitmq"
	"github.com/samarthasthan/tanx-task/pkg/env"
	rabbitmq_utils "github.com/samarthasthan/tanx-task/pkg/rabbitmq"
)

// Define the environment variables
var (
	MYSQL_PORT            string
	MYSQL_ROOT_PASSWORD   string
	MYSQL_HOST            string
	REDIS_PORT            string
	REDIS_HOST            string
	RABBITMQ_DEFAULT_USER string
	RABBITMQ_DEFAULT_PASS string
	RABBITMQ_DEFAULT_PORT string
	RABBITMQ_DEFAULT_HOST string
)

// Initialize the environment variables
func init() {
	MYSQL_PORT = env.GetEnv("MYSQL_PORT", "3306")
	MYSQL_ROOT_PASSWORD = env.GetEnv("MYSQL_ROOT_PASSWORD", "password")
	MYSQL_HOST = env.GetEnv("MYSQL_HOST", "localhost")
	REDIS_PORT = env.GetEnv("REDIS_PORT", "6379")
	REDIS_HOST = env.GetEnv("REDIS_HOST", "localhost")
	RABBITMQ_DEFAULT_USER = env.GetEnv("RABBITMQ_DEFAULT_USER", "root")
	RABBITMQ_DEFAULT_PASS = env.GetEnv("RABBITMQ_DEFAULT_PASS", "password")
	RABBITMQ_DEFAULT_PORT = env.GetEnv("RABBITMQ_DEFAULT_PORT", "5672")
	RABBITMQ_DEFAULT_HOST = env.GetEnv("RABBITMQ_DEFAULT_HOST", "localhost")
}

func main() {

	// Define the combined WebSocket URL for multiple trading pairs
	// Replace the pairs with the ones you're interested in
	streams := "btcusdt@trade"
	combinedURL := url.URL{
		Scheme:   "wss",
		Host:     "stream.binance.com:9443",
		Path:     "/stream",
		RawQuery: fmt.Sprintf("streams=%s", streams),
	}

	// Connect to the WebSocket
	fmt.Println("Connecting to", combinedURL.String())
	conn, _, err := websocket.DefaultDialer.Dial(combinedURL.String(), nil)
	if err != nil {
		log.Fatal("Failed to connect to WebSocket:", err)
	}
	defer conn.Close()
	fmt.Println("Connected to WebSocket")

	conn.SetReadDeadline(time.Time{})

	// Create a new instance of MySQL database
	mysql := database.NewMySQL()
	err = mysql.Connect(fmt.Sprintf("root:%s@tcp(%s:%s)/tanx?parseTime=true", MYSQL_ROOT_PASSWORD, MYSQL_HOST, MYSQL_PORT))
	if err != nil {
		panic(err.Error())
	}
	defer mysql.Close()

	// Create a new instance of Redis database
	redis := database.NewRedis()
	err = redis.Connect(fmt.Sprintf("%s:%s", REDIS_HOST, REDIS_PORT))
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

	// Alertor is a struct which contains the instances of Redis and RabbitMQ
	// Create a new instance of Alertor
	// Pass the instance of RabbitMQ and Redis to the Alertor
	// This will allow the Alertor to interact with Redis and RabbitMQ
	// The Alertor will be responsible for sending alerts to the RabbitMQ
	alertor := alert.NewAlert(publisher, redis, mysql, conn)

	// Set up a WebSocket reader
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				break
			}
			alertor.ProcessMessage(message)
		}
	}()

	// Keep the main function running
	select {}
}
