package webmodel

import (
	"testing"
)

func TestValidate(t *testing.T) {
	UCs := []UserCredentials{
		{
			UserName:    "usserAdd",
			Email:       "emailtAdd@email",
			Password:    "pass",
			DateOfBirth: "2002-03-03",
			Gender:      "He",
			FirstName:   "John",
			LastName:    "second",
		},
		{
			UserName:    "",
			Email:       "emailtAdd@email",
			Password:    "pass",
			DateOfBirth: "2002-03-03",
			Gender:      "He",
			FirstName:   "John",
			LastName:    "second",
		},
		{
			UserName:    "usserAdd",
			Email:       "",
			Password:    "pass",
			DateOfBirth: "2002-03-03",
			Gender:      "He",
			FirstName:   "John",
			LastName:    "second",
		},
		{
			UserName:    "usserAdd",
			Email:       "emailtAdd@email",
			Password:    "",
			DateOfBirth: "2002-03-03",
			Gender:      "He",
			FirstName:   "John",
			LastName:    "second",
		},
		{
			UserName:    "usserAdd",
			Email:       "emailtAdd@email",
			Password:    "pass",
			DateOfBirth: "",
			Gender:      "He",
			FirstName:   "John",
			LastName:    "second",
		},
		{
			UserName:    "usserAdd",
			Email:       "emailtAdd@email",
			Password:    "pass",
			DateOfBirth: "2002-03-03",
			Gender:      "",
			FirstName:   "John",
			LastName:    "second",
		},
		{
			UserName:    "usserAdd",
			Email:       "emailtAdd@email",
			Password:    "pass",
			DateOfBirth: "2002-03-03",
			Gender:      "He",
			FirstName:   "",
			LastName:    "second",
		},
		{
			UserName:    "usserAdd",
			Email:       "emailtAdd@email",
			Password:    "pass",
			DateOfBirth: "2002-03-03",
			Gender:      "He",
			FirstName:   "John",
			LastName:    "",
		},
		{
			UserName:    "usserAdd",
			Email:       "emailtAddemail",
			Password:    "pass",
			DateOfBirth: "2002-03-03",
			Gender:      "He",
			FirstName:   "John",
			LastName:    "second",
		},
	}

	message := []string{
		"",
		"username missing",
		"email missing",
		"password missing",
		"dateBirth missing",
		"gender missing",
		"firstName missing",
		"lastName missing",
		"wrong email",
	}

	for i := 0; i < len(UCs); i++ {
		res := UCs[i].Validate()
		if res != message[i] {
			t.Fatalf("# %d: result is\n%s\n, but expected to be\n%s\n", i, res, message[i])
		}
	}
}
