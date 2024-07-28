package alertor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/samarthasthan/tanx-task/internal/database"
	"github.com/samarthasthan/tanx-task/internal/database/mysql/sqlc"
	"github.com/samarthasthan/tanx-task/internal/models"
	"github.com/samarthasthan/tanx-task/internal/rabbitmq"
)

type Alert struct {
	redis    database.Database
	mysql    database.Database
	rabbitmq *rabbitmq.RabbitMQ
	ws       *websocket.Conn
}

func NewAlert(rb *rabbitmq.RabbitMQ, redis database.Database, mysql database.Database, ws *websocket.Conn) *Alert {
	return &Alert{rabbitmq: rb, redis: redis, mysql: mysql, ws: ws}
}

// Process incoming WebSocket messages
func (a *Alert) ProcessMessage(message []byte) {
	mysql := a.mysql.(*database.MySQL)
	redis := a.redis.(*database.Redis)
	dbCtx := context.Background()

	// Parse the message to extract relevant information
	var msg map[string]interface{}
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Println("Error parsing message:", err)
		return
	}

	data, ok := msg["data"].(map[string]interface{})
	if !ok {
		log.Println("Invalid message format: missing 'data' field")
		return
	}

	symbol, okSymbol := data["s"].(string)
	price, okPrice := data["p"].(string)
	if !okSymbol || !okPrice {
		log.Println("Invalid message format: missing 's' or 'p' field")
		return
	}

	// Define the cache key for Redis
	cacheKey := "alerts"

	// Attempt to fetch alerts from Redis
	cachedAlerts, err := redis.Get(dbCtx, cacheKey).Result()
	if err == nil && cachedAlerts != "" {
		// Cache hit: Unmarshal cached alerts
		var alerts []sqlc.GetAlertsRow // Adjust the type according to your actual data model
		if err := json.Unmarshal([]byte(cachedAlerts), &alerts); err == nil {
			a.processAlerts(alerts, symbol, price)
			return
		}
	}

	// Cache miss or error: Fetch alerts from MySQL
	alerts, err := mysql.Queries.GetAlerts(dbCtx)
	if err != nil {
		log.Println("Error fetching alerts from MySQL:", err)
		return
	}

	alertsJSON, err := json.Marshal(alerts)
	if err != nil {
		log.Println("Error marshaling alerts:", err)
		return
	}

	// Cache the fetched alerts in Redis with an expiration time
	err = redis.Set(dbCtx, cacheKey, alertsJSON, 5*time.Minute).Err()
	if err != nil {
		log.Println("Error caching alerts in Redis:", err)
	}

	// Process the alerts
	a.processAlerts(alerts, symbol, price)
}

func (a *Alert) processAlerts(alerts []sqlc.GetAlertsRow, symbol, priceStr string) {
	mysql := a.mysql.(*database.MySQL)
	redis := a.redis.(*database.Redis)

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		log.Println("Error converting price to float:", err)
		return
	}

	// Define an absolute tolerance value
	const toleranceAmount = 90.0 // Absolute tolerance of 100 units

	for _, alert := range alerts {
		if symbol == alert.Curreny {
			alertPrice, err := strconv.ParseFloat(alert.Price, 64)
			if err != nil {
				log.Println("Error converting alert price to float:", err)
				continue
			}

			// Check if the real price is within the tolerance range of the alert price
			if math.Abs(price-alertPrice) <= toleranceAmount {
				// Create a new Mail struct
				mail := &models.Mail{
					To:      alert.Email,
					Subject: "Trigger Alert",
					Body:    fmt.Sprintf("Alert triggered for %s at %s", alert.Curreny, alert.Price),
				}

				// struct to byte
				data, err := json.Marshal(mail)
				if err != nil {
					log.Println(err.Error())
				}

				// Publish a message to the RabbitMQ
				err = a.rabbitmq.Publish("tanx", "tanx", data)
				if err != nil {
					log.Println(err.Error())
				}

				// Update the alert status in MySQL
				dbCtx := context.Background()
				err = mysql.Queries.UpdateAlertStatus(dbCtx, sqlc.UpdateAlertStatusParams{
					Alertid: alert.Alertid,
					Status:  "triggered",
				})
				if err != nil {
					log.Println("Error updating alert status in MySQL:", err)
				}

				// Remove cached alerts from Redis
				err = redis.Del(dbCtx, "alerts").Err()
				if err != nil {
					log.Printf("Error removing alerts cache from Redis: %v", err)
				}

			}
		}
	}
}
