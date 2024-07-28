package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samarthasthan/tanx-task/internal/models"
)

// Alert alerts/create/, alerts/delete/, alerts/all/

// AlertCreate handles the alert creation request
func (h *Handlers) handleAlertCreate(c echo.Context) error {
	s := new(models.Alert)
	if err := c.Bind(s); err != nil {
		return err
	}
	if err := c.Validate(s); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := h.controller.CreateAlert(c, s); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(200, map[string]string{"message": "Alert created successfully"})
}

// AlertDelete handles the alert deletion request
func (h *Handlers) handleAlertDelete(c echo.Context) error {
	s := new(models.DeleteAlert)
	if err := c.Bind(s); err != nil {
		return err
	}
	if err := c.Validate(s); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := h.controller.DeleteAlert(c, s); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(200, map[string]string{"message": "Alert deleted successfully"})
}

// AlertAll handles the alert all request
func (h *Handlers) handleAlertAll(c echo.Context) error {
	alerts, err := h.controller.GetAlerts(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(200, alerts)
}
