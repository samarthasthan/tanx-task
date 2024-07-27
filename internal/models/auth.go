package models

type SignUp struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type VerifyOTP struct {
	Email string `json:"email" validate:"required,email"`
	OTP   int    `json:"otp" validate:"required"`
}
