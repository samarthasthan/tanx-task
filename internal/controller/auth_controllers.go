package controller

import (
	"time"

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

	err = mysql.Queries.CreateAccount(dbCtx, sqlc.CreateAccountParams{
		Userid:   uuid.New().String(),
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
		Userid:         u.Email,
		Otp:            int32(OTP),
		Expiresat:      time.Now().Add(OTP_EXPIRATION_TIME),
	})

	if err != nil {
		return err
	}

	return nil
}
