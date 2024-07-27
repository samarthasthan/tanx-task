package api

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/samarthasthan/tanx-task/internal/controller"
)

type (
	CustomValidator struct {
		validator *validator.Validate
	}
	Handlers struct {
		*echo.Echo
		controller *controller.Controller
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func NewHandler(c *controller.Controller) *Handlers {
	return &Handlers{Echo: echo.New(), controller: c}
}

func (h *Handlers) RegisterValidator() {
	h.Validator = &CustomValidator{validator: validator.New()}
}

func (h *Handlers) Handle() {
	// Handle the authentication routes
	h.RegisterValidator()
	h.HandleAuth()
}

func (h *Handlers) HandleAuth() {
	auth := h.Group("/auth")
	auth.POST("/signup", h.handleSignUp)
	auth.POST("/otp-verification", h.handleOTPVerification)
	auth.POST("/login", h.handleLogin)

	authRes := h.Group("/auth")
	authRes.Use(h.validateToken)
	authRes.POST("/verify", h.handleVerify)
}
