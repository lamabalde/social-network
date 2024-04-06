package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Notification struct {
	ID         int            `json:"id"`
	UserID     string         `json:"userId"`
	FromUserID sql.NullString `json:"fromUserId,omitempty"`
	Type       string         `json:"type"`
	Body       string         `json:"body"`
	GroupID    sql.NullString `json:"groupId,omitempty"`
	PostID     sql.NullString `json:"postId,omitempty"`
	Read       int            `json:"read"`
	DateCreate time.Time      `json:"dateCreate"`
}

func (n *Notification) String() string {
	if n == nil {
		return "nil"
	}

	return fmt.Sprintf("\nNotification: to user id '%s' type '%s'\n    text: '%s'\n    FromUserID: '%v'\n    GroupID: '%v'\n    PostID:'%v'\n    Created: %s\n isread: %d\n",
		n.UserID, n.Type, n.Body, n.FromUserID, n.GroupID, n.PostID, n.DateCreate, n.Read)
}
