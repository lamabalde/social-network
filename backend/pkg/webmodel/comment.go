package webmodel

import (
	"time"
)

type Comment struct {
	PostID     string    `json:"postID"`
	Content    string    `json:"content"`
	DateCreate time.Time `json:"date,omitempty"`
	Image      string    `json:"image,omitempty"`
}

type CommentImg struct {
	Image     string
	CommentID string
}

func (c *Comment) Validate() string {
	if IsEmpty(c.Content) {
		return "Comment's text missing"
	}

	// if c.DateCreate.Before(time.Date(2023, time.September, 1, 0, 0, 0, 0, time.UTC)) {
	// 	return "Date is too old"
	// }
	return ""
}
