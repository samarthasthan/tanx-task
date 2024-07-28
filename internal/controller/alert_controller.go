package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samarthasthan/tanx-task/internal/database"
	"github.com/samarthasthan/tanx-task/internal/database/mysql/sqlc"
	"github.com/samarthasthan/tanx-task/internal/models"
)

func (c *Controller) CreateAlert(ctx echo.Context, a *models.Alert) error {
	mysql := c.mysql.(*database.MySQL)
	redis := c.redis.(*database.Redis)
	dbCtx := ctx.Request().Context()
	userID := ctx.Get("id").(string)

	err := mysql.Queries.CreateAlert(dbCtx, sqlc.CreateAlertParams{
		Userid:  userID,
		Alertid: uuid.New().String(),
		Curreny: a.Currency,
		Price:   fmt.Sprintf("%f", a.Price),
	})

	if err != nil {
		return err
	}

	// Remove cached alerts from Redis
	err = redis.Del(dbCtx, "alerts").Err()
	if err != nil {
		log.Printf("Error removing alerts cache from Redis: %v", err)
		// Optionally, handle this error depending on your use case
	} else {
		log.Printf("Alert created successfully for user %s and cache invalidated", userID)
	}

	return nil
}

func (c *Controller) DeleteAlert(ctx echo.Context, a *models.DeleteAlert) error {
	mysql := c.mysql.(*database.MySQL)
	dbCtx := ctx.Request().Context()

	err := mysql.Queries.DeleteAlert(dbCtx, a.AlertID)

	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) GetAlerts(ctx echo.Context) ([]sqlc.GetAlertsByUserRow, error) {
	mysql := c.mysql.(*database.MySQL)
	redis := c.redis.(*database.Redis)
	dbCtx := ctx.Request().Context()
	userID := ctx.Get("id").(string)

	// Define the cache key
	cacheKey := "alerts:" + userID

	cachedAlerts, err := redis.Get(dbCtx, cacheKey).Result()
	if err == nil && cachedAlerts != "" {
		// If cache hit, unmarshal the cached data
		var alerts []sqlc.GetAlertsByUserRow
		if err := json.Unmarshal([]byte(cachedAlerts), &alerts); err == nil {
			return alerts, nil
		}
	}

	alerts, err := mysql.Queries.GetAlertsByUser(dbCtx, userID)
	if err != nil {
		return nil, err
	}

	alertsJSON, err := json.Marshal(alerts)
	if err != nil {

		return nil, err
	}

	// Store the alerts in Redis cache with an expiration time
	err = redis.Set(dbCtx, cacheKey, alertsJSON, 5*time.Minute).Err()
	if err != nil {
		return nil, err
	}

	return alerts, nil
}
