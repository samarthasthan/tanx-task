package mail

import (
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/labstack/gommon/log"
	"github.com/samarthasthan/tanx-task/internal/models"
	"github.com/samarthasthan/tanx-task/internal/rabbitmq"
	"gopkg.in/gomail.v2"
)

type MailHandler struct {
	rabbitmq *rabbitmq.RabbitMQ

	SMTP_SERVER   string
	SMTP_PORT     string
	SMTP_LOGIN    string
	SMTP_PASSWORD string
}

func NewMailHandler(
	consumer *rabbitmq.RabbitMQ,
	SMTP_SERVER string,
	SMTP_PORT string,
	SMTP_LOGIN string,
	SMTP_PASSWORD string,
) *MailHandler {
	return &MailHandler{
		rabbitmq:      consumer,
		SMTP_SERVER:   SMTP_SERVER,
		SMTP_PORT:     SMTP_PORT,
		SMTP_LOGIN:    SMTP_LOGIN,
		SMTP_PASSWORD: SMTP_PASSWORD,
	}
}

func (h *MailHandler) SendMail(in *models.Mail) error {
	from := h.SMTP_LOGIN                   // Sender email
	to := in.To                            // Recipient email
	host := h.SMTP_SERVER                  // SMTP server
	port, err := strconv.Atoi(h.SMTP_PORT) // SMTP port

	if err != nil {
		return err
	}

	// Message
	msg := gomail.NewMessage()
	msg.SetHeader("From", "noreply@samarthasthan.com") // Sender email
	msg.SetHeader("To", to)                            // Recipient email
	msg.SetHeader("Subject", in.Subject)               // Subject of the email
	// text/html for a html email
	msg.SetBody("text/html", in.Body) // Body of the email

	n := gomail.NewDialer(host, port, from, h.SMTP_PASSWORD) // SMTP server details

	// Send the email
	if err := n.DialAndSend(msg); err != nil {
		return err
	}

	log.Infof("Mail sent to %s", to)
	return nil
}

func (h *MailHandler) ConsumeMails() {
	for {
		// Consume a message from the RabbitMQ
		msgs, err := h.rabbitmq.Consume("tanx", "tanx")
		if err != nil {
			log.Errorf("Failed to consume a message: %s", err)
		}
		for d := range msgs {
			mail := models.Mail{}
			bytes := bytes.NewBuffer(d.Body)
			json.NewDecoder(bytes).Decode(&mail)
			if err := h.SendMail(&mail); err != nil {
				log.Errorf("Failed to send mail: %s", err)
			}
		}
	}
}
