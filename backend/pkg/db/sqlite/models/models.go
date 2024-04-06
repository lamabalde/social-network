package models

import (
	"errors"
	"fmt"
)

const (
	POST    = "post"
	COMMENT = "comment"
)

const N_REACTIONS = 2

const (
	DISLIKE UserReactions = iota
	LIKE
)

const (
	CHAT_TYPE_PRIVATE = iota
	CHAT_TYPE_GROUP
)

var (
	ErrNoRecords       = errors.New("there is no record in the DB")
	ErrTooManyRecords  = errors.New("there are more than one record")
	ErrUnique          = errors.New("unique constraint failed")
	ErrUniqueUserName  = errors.New("user with the given name already exists")
	ErrUniqueUserEmail = errors.New("user with the given email already exists")
	ErrAddConstaints   = errors.New(" ADD to table constraint")
)

type ErrorMessage struct {
	Errors string
}

type UserReactions int

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

func (c *Category) String() string {
	if c == nil {
		return "nil"
	}
	return fmt.Sprintf("categ id: %s | name: %s\n", c.ID, c.Name)
}

// type GroupMessage struct {
// 	ID         string    `json:"id"`
// 	Author     *User     `json:"author,omitempty"`
// 	Content    string    `json:"content"`
// 	DateCreate time.Time `json:"dateCreate,omitempty"`
// 	Images     []string  `json:"-"`
// }

/*
	type Message struct {
		MessageID   string
		Body        string `json:"body"`
		SenderID    string `json:"sender"`
		ReceiverID  string `json:"Receiver"`
		MessageType string `json:"messageType"` // msg, login, logout
	}

	type Connection struct {
		client   websocket.Conn
		userName string
		userID   string
	}
*/
type Migration struct {
	Version int
	Dirty   int
}

/*
checkID must return true if a user with the given ID has to be added to the result of selection from DB
*/
type IdChecker interface {
	CheckID(id string) bool
}
