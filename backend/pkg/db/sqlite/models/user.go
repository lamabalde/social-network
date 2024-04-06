package models

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	FOLLOW_STATUS_FOLLOWING = "following"
	FOLLOW_STATUS_REQUESTED = "requested"
)

type User struct {
	ID              string     `json:"id"`
	UserName        string     `json:"userName"`
	Password        string     `json:"password,omitempty"`
	Email           string     `json:"email,omitempty"`
	DateCreate      time.Time  `json:"dateCreate,omitempty"`
	DateBirth       time.Time  `json:"dateOfBirth,omitempty"`
	Gender          string     `json:"gender,omitempty"`
	FirstName       string     `json:"firstName,omitempty"`
	LastName        string     `json:"lastName,omitempty"`
	ProfileType     int        `json:"profileType"` // 0 - public, 1 - private
	Followers       []UserBase `json:"followers,omitempty"`
	Followings      []UserBase `json:"followings,omitempty"`
	LastMessageDate string     `json:"lastMessageDate,omitempty"`
	ProfileImg      string     `json:"profileImg,omitempty"`
	AboutMe         string     `json:"aboutMe,omitempty"`
}

type UserBase struct {
	ID       string `json:"id"`
	UserName string `json:"userName"`
}

func (u *User) String() string {
	if u == nil {
		return "nil"
	}
	return fmt.Sprintf(`name: %s (id: %s)`, u.UserName, u.ID)
}

func (u *User) StringFull() string {
	if u == nil {
		return "nil"
	}
	return fmt.Sprintf(`user: {
		    id:           %s
		    username:     %s
		    email:        %s
		    password:     %s
		    DataCreate:   %s
		    DateBirth:    %s
		    Gender:       %s
		    First Name:   %s
		    Last Name:    %s
		    ProfileType	: %d
		    Followers:    %s
		    Followings:   %s
			}`,
		u.ID, u.UserName, u.Email, "****", u.DateCreate.String(), u.DateBirth, u.Gender, u.FirstName, u.LastName, u.ProfileType, u.Followers, u.Followings)
}

func (ub *UserBase) String() string {
	return fmt.Sprintf(`%s (id: %s)`, ub.UserName, ub.ID)
}

func (u *User) UnmarshalJSON(data []byte) error {
	type tmpUser User
	var tmp tmpUser
	err := json.Unmarshal(data, &tmp)
	*u = User(tmp)
	return err
}
