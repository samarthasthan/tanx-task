package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samarthasthan/tanx-task/internal/database"
	"github.com/samarthasthan/tanx-task/internal/database/mysql/sqlc"
	"github.com/samarthasthan/tanx-task/internal/models"
	bcrpyt "github.com/samarthasthan/tanx-task/pkg/bycrpt"
	"github.com/samarthasthan/tanx-task/pkg/otp"
)

func (c *Controller) SignUp(ctx echo.Context, u *models.SignUp) error {
	mysql := c.mysql.(*database.MySQL)

	// Hash the password
	hashedPassword, err := bcrpyt.HashPassword(u.Password)
	if err != nil {
		return err
	}

	dbCtx := ctx.Request().Context()

	userID := uuid.New().String()

	err = mysql.Queries.CreateAccount(dbCtx, sqlc.CreateAccountParams{
		Userid:   userID,
		Name:     u.Name,
		Email:    u.Email,
		Password: hashedPassword,
	})

	if err != nil {
		return err
	}

	OTP := otp.GenerateVerificationCode()

	err = mysql.Queries.CreateVerification(dbCtx, sqlc.CreateVerificationParams{
		Verificationid: uuid.New().String(),
		Userid:         userID,
		Otp:            int32(OTP),
		Expiresat:      time.Now().Add(OTP_EXPIRATION_TIME),
	})

	if err != nil {
		return err
	}

	// Create a new Mail struct
	mail := &models.Mail{
		To:      u.Email,
		Subject: "Welcome to TanX",
		Body:    fmt.Sprintf("Your OTP for TanX registration is %d", OTP),
	}

	// struct to byte
	data, err := json.Marshal(mail)
	if err != nil {
		log.Println(err.Error())
	}

	// Publish a message to the RabbitMQ
	err = c.rabbitmq.Publish("tanx", "tanx", data)
	if err != nil {
		log.Println(err.Error())
	}

	return nil
}

func (c *Controller) VerifyOTP(ctx echo.Context, v *models.VerifyOTP) error {
	mysql, ok := c.mysql.(*database.MySQL)
	if !ok {
		return fmt.Errorf("mysql is not of type *database.MySQL")
	}

	dbCtx := ctx.Request().Context()
	tx, err := mysql.DB.BeginTx(dbCtx, nil)

	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	// Get UserID from Email
	userID, err := mysql.Queries.GetUserIDByEmail(dbCtx, v.Email)
	if err != nil {
		return fmt.Errorf("failed to get user id: %v", err)
	}

	// Get OTP from database
	otpRow, err := mysql.Queries.GetOTP(dbCtx, userID)
	if err != nil {
		return fmt.Errorf("failed to get OTP: %v", err)
	}

	// Check if OTP is valid
	if otpRow.Otp != int32(v.OTP) {
		return fmt.Errorf("invalid OTP")
	}

	// Check if OTP has expired
	if otp.CheckOTPExpiration(otpRow.Expiresat) == false {
		return fmt.Errorf("OTP has expired")
	}

	// verify user
	err = mysql.Queries.VerifyAccount(dbCtx, userID)
	if err != nil {
		return fmt.Errorf("failed to verify account: %v", err)
	}

	err = mysql.Queries.DeleteVerification(dbCtx, userID)

	return nil
}

func (c *Controller) Login(ctx echo.Context, l *models.Login) (string, error) {
	// Get the user's password from the database
	mysql := c.mysql.(*database.MySQL)

	dbCtx := ctx.Request().Context()

	user, err := mysql.Queries.GetUserByEmail(dbCtx, l.Email)
	if err != nil {
		return "", err
	}

	// Compare the user's password with the provided password
	if !bcrpyt.ValidatePassword(user.Password, l.Password) {
		return "", fmt.Errorf("invalid password")
	}

	// Create a JWT token
	token, err := c.createToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Helper function to create a JWT token
func (c *Controller) createToken(u sqlc.GetUserByEmailRow) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":         u.Userid,
			"name":       u.Name,
			"email":      u.Email,
			"expires_at": time.Now().Add(356 * time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString([]byte(c.jwt_secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Validate the JWT token
func (c *Controller) VerifyToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(c.jwt_secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
