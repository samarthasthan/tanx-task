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

type Controller struct {
	mysql    database.Database
	rabbitmq *rabbitmq.RabbitMQ
}

func NewController(rb *rabbitmq.RabbitMQ, db database.Database) *Controller {
	return &Controller{rabbitmq: rb, mysql: db}
}
