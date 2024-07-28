package controller

import (
	"fmt"

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
	dbCtx := ctx.Request().Context()
	userID := ctx.Get("id").(string)

	alerts, err := mysql.Queries.GetAlerts(dbCtx, userID)

	if err != nil {
		return nil, err
	}

	return alerts, nil
}
