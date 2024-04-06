package models

import (
	"fmt"
	"time"
)

type Event struct {
	ID          string               `json:"id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	DateCreate  time.Time            `json:"dateCreate"`
	DateEvent   time.Time            `json:"dateEvent"`
	GroupID     string               `json:"groupID"`
	CreatorID   string               `json:"creatorID"`
	CreatorName string               `json:"creatorName"`
	UserOptions []UserOptionForEvent `json:"userOptions"`
}

type UserOptionForEvent struct {
	UserID   string `json:"userID"`
	UserName string `json:"userName"`
	Option   int    `json:"option"`
}

func (e *Event) String() string {
	if e == nil {
		return "nil"
	}

	return fmt.Sprintf("\nEvent: '%s' (id: '%s') -----------\n    Description:'%s'\n    Date:%s\n    Created:%s\n    in group id:'%s'\n    by user id:'%s'",
		e.Title, e.ID, e.Description, e.DateEvent, e.DateCreate, e.GroupID, e.CreatorID)
}
