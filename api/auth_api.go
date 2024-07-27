package api

import "github.com/labstack/echo/v4"

// SignUp handles the sign up request
func (h *API) SignUp(c echo.Context) error {
	return c.String(200, "Sign Up")
}

// OTPVerification handles the OTP verification request
func (h *API) OTPVerification(c echo.Context) error {
	return c.String(200, "OTP Verification")
}

// Login handles the login request
func (h *API) Login(c echo.Context) error {
	return c.String(200, "Login")
}
