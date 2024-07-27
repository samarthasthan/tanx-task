package main

import (
	"fmt"
	"log"

	"github.com/samarthasthan/tanx-task/api"
	"github.com/samarthasthan/tanx-task/internal/database"
	"github.com/samarthasthan/tanx-task/pkg/env"
)

// Define the environment variables
var (
	REST_API_PORT       string
	MYSQL_PORT          string
	MYSQL_ROOT_PASSWORD string
	MYSQL_HOST          string
	REDIS_PORT          string
	REDIS_HOST          string
	KAFKA_PORT          string
	KAFKA_EXTERNAL_PORT string
	KAFKA_HOST          string
)

// Initialize the environment variables
func init() {
	REST_API_PORT = env.GetEnv("REST_API_PORT", "3000")
	MYSQL_PORT = env.GetEnv("MYSQL_PORT", "3306")
	MYSQL_ROOT_PASSWORD = env.GetEnv("MYSQL_ROOT_PASSWORD", "root")
	MYSQL_HOST = env.GetEnv("MYSQL_HOST", "mysql")
	REDIS_PORT = env.GetEnv("REDIS_PORT", "6379")
	REDIS_HOST = env.GetEnv("REDIS_HOST", "redis")
	KAFKA_PORT = env.GetEnv("KAFKA_PORT", "9092")
	KAFKA_EXTERNAL_PORT = env.GetEnv("KAFKA_EXTERNAL_PORT", "9094")
	KAFKA_HOST = env.GetEnv("KAFKA_HOST", "kafka")
}

func main() {
	// Create a new instance of MySQL database
	mysql := database.NewMySQL()
	err := mysql.Connect(fmt.Sprintf("root:%s@tcp(%s:%s)/authentication?parseTime=true", MYSQL_ROOT_PASSWORD, MYSQL_HOST, MYSQL_PORT))
	if err != nil {
		log.Println(err.Error())
	}

	// Register Echo handler
	h := api.NewAPI()
	h.Handle()
	h.Logger.Fatal(h.Start(":" + REST_API_PORT))
}
