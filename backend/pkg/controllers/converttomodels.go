package controllers

import (
	"fmt"
	"time"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/helpers"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel"
)

func CreateUser(u webmodel.UserCredentials) (*models.User, error) {
	user := models.User{
		UserName:   u.UserName,
		Email:      u.Email,
		DateCreate: time.Now(),
		Gender:     u.Gender,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		AboutMe:    u.AboutMe,
	}
	var err error

	user.Password, err = helpers.HashPassword(u.Password) // hash the password
	if err != nil {
		return nil, fmt.Errorf("failed to generate crypto password: %w", err)
	}

	user.DateBirth, err = time.Parse(time.DateOnly, u.DateOfBirth) // date must by in the format  "2006-01-02"
	if err != nil {
		return nil, fmt.Errorf("failed to parse the date of birth: %w", err)
	}

	return &user, nil
}
