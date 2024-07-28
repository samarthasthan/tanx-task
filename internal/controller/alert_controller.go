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

func (c *Controller) GetAllAlerts(ctx echo.Context) ([]sqlc.GetAlertsRow, error) {
	mysql := c.mysql.(*database.MySQL)
	redis := c.redis.(*database.Redis)
	dbCtx := ctx.Request().Context()
	userID := ctx.Get("id").(string)

	// Define the cache key
	cacheKey := "alerts:" + userID

	cachedAlerts, err := redis.Get(dbCtx, cacheKey).Result()
	if err == nil && cachedAlerts != "" {
		// If cache hit, unmarshal the cached data
		log.Println("Cache hit")
		var alerts []sqlc.GetAlertsRow
		if err := json.Unmarshal([]byte(cachedAlerts), &alerts); err == nil {
			return alerts, nil
		}
	}

	alerts, err := mysql.Queries.GetAlerts(dbCtx, userID)
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
