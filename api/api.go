package api

import (
	"github.com/labstack/echo/v4"
)

type API struct {
	*echo.Echo
}

func NewAPI() *API {
	return &API{echo.New()}
}

func (h *API) Handle() {
	// Handle the authentication routes
	h.HandleAuth()
}

func (h *API) HandleAuth() {
	auth := h.Group("/auth")
	auth.POST("/signup", h.SignUp)
	auth.POST("/otp-verification", h.OTPVerification)
	auth.POST("/login", h.Login)
}
