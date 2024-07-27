package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samarthasthan/tanx-task/internal/models"
)

// SignUp handles the sign up request
func (h *Handlers) handleSignUp(c echo.Context) error {
	s := new(models.SignUp)
	if err := c.Bind(s); err != nil {
		return err
	}
	if err := c.Validate(s); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := h.controller.SignUp(c, s); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(200, map[string]string{"message": "Account created successfully"})
}

// OTPVerification handles the OTP verification request
func (h *Handlers) handleOTPVerification(c echo.Context) error {
	return c.String(200, "OTP Verification")
}

// Login handles the login request
func (h *Handlers) handleLogin(c echo.Context) error {
	return c.String(200, "Login")
}
