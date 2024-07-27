package api

import (
	"net/http"
	"strings"

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
	s := new(models.VerifyOTP)
	if err := c.Bind(s); err != nil {
		return err
	}
	if err := c.Validate(s); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := h.controller.VerifyOTP(c, s); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(200, map[string]string{"message": "OTP verified successfully"})
}

// Login handles the login request
func (h *Handlers) handleLogin(c echo.Context) error {
	l := new(models.Login)
	if err := c.Bind(l); err != nil {
		return err
	}
	if err := c.Validate(l); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	token, err := h.controller.Login(c, l)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func (h *Handlers) handleVerify(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"id":    c.Get("id"),
		"name":  c.Get("name"),
		"email": c.Get("email"),
	})
}

// Echo middleware to validate the JWT token
func (h *Handlers) validateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"error": "missing token",
			})
		}

		// Remove the "Bearer " prefix
		token = strings.TrimPrefix(token, "Bearer ")

		claims, err := h.controller.VerifyToken(token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"error": err.Error(),
			})
		}

		c.Set("id", claims["id"])
		c.Set("name", claims["name"])
		c.Set("email", claims["email"])
		return next(c)
	}
}
