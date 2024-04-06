package controllers

import (
	"testing"
	"time"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel"
)

func TestCreateUser(t *testing.T) {
	userExpected := models.User{
		UserName:   "usserAdd",
		Email:      "emailtAdd@email",
		DateCreate: time.Now(),
		DateBirth:  time.Date(2002, time.March, 3, 0, 0, 0, 0, time.UTC),
		Gender:     "He",
		FirstName:  "John",
		LastName:   "second",
	}

	uc := webmodel.UserCredentials{
		UserName:    "usserAdd",
		Email:       "emailtAdd@email",
		DateOfBirth: "2002-03-03",
		Gender:      "He",
		FirstName:   "John",
		LastName:    "second",
	}
	user, err := CreateUser(uc)
	if err != nil {
		t.Fatalf("err is %s\n", err)
	}
	if userExpected.UserName != user.UserName || userExpected.Email != user.Email || userExpected.DateBirth != user.DateBirth || userExpected.Gender != user.Gender || userExpected.FirstName != user.FirstName || userExpected.LastName != user.LastName {
		t.Fatalf("result is\n%s\n, but expected to be\n%s\n", user, &userExpected)
	}
}
