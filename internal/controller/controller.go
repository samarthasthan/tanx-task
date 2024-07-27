package controller

import (
	"time"

	"github.com/samarthasthan/tanx-task/internal/database"
)

var (
	OTP_EXPIRATION_TIME time.Duration
)

func init() {
	OTP_EXPIRATION_TIME = 5 * time.Minute
}

type Controller struct {
	mysql database.Database
}

func NewController(db database.Database) *Controller {
	return &Controller{mysql: db}
}
