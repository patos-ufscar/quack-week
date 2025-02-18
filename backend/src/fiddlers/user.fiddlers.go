package fiddlers

import (
	"time"

	"github.com/patos-ufscar/quack-week/common"
	"github.com/patos-ufscar/quack-week/models"
)

func NewUnconfirmedUser(email string, password string, firstName string, lastName string, dateOfBirth *time.Time) (*models.UnconfirmedUser, error) {
	hash, err := common.HashPassword(password)
	if err != nil {
		return nil, err
	}

	otp, err := common.GenerateRandomString(common.OTP_LEN)
	if err != nil {
		return nil, err
	}

	return &models.UnconfirmedUser{
		Email:        email,
		Otp:          otp,
		PasswordHash: hash,
		FirstName:    firstName,
		LastName:     lastName,
		DateOfBirth:  dateOfBirth,
	}, nil
}
