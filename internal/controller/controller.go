package controller

import (
	"time"

	"github.com/samarthasthan/tanx-task/internal/database"
	"github.com/samarthasthan/tanx-task/internal/rabbitmq"
)

var (
	OTP_EXPIRATION_TIME time.Duration
)

func init() {
	OTP_EXPIRATION_TIME = 5 * time.Minute
}

// For Dependency Injection we are using database.Database interface instead of concrete type
type Controller struct {
	mysql      database.Database
	redis      database.Database
	rabbitmq   *rabbitmq.RabbitMQ
	jwt_secret string
}

func NewController(rb *rabbitmq.RabbitMQ, mysql database.Database, redis database.Database, jwt string) *Controller {
	return &Controller{rabbitmq: rb, mysql: mysql, redis: redis, jwt_secret: jwt}
}
