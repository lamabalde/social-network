package webmodel

import (
	"time"
)

type Post struct {
	Theme      string    `json:"title"`
	Content    string    `json:"body"`
	Categories []string  `json:"categories,omitempty"`
	DateCreate time.Time `json:"date,omitempty"`
	GroupID    string    `json:"groupID,omitempty"`
	Image      string    `json:"postImage,omitempty"`
	PostType   int       `json:"postType,omitempty"`
}

func (p *Post) Validate() string {
	if IsEmpty(p.Theme) {
		return "Post's theme missing"
	}
	if IsEmpty(p.Content) {
		return "Post's text missing"
	}

	// if p.DateCreate.Before(time.Date(2023, time.September, 1, 0, 0, 0, 0, time.UTC)) {
	// 	return "Date is too old"
	// }

	return ""
}

type PostsInGroupPortion struct {
	GroupID string `json:"groupID"`
	Offset  int    `json:"offset"`
}

func (p *PostsInGroupPortion) Validate() string {
	if IsEmpty(p.GroupID) {
		return "GroupID missing"
	}
	if p.Offset < 0 {
		return "negative offset"
	}

	return ""
}

type UserPosts struct {
	UserID string `json:"userID"`
	Offset int    `json:"offset"`
}

func (p *UserPosts) Validate() string {
	if IsEmpty(p.UserID) {
		return "UserID missing"
	}
	if p.Offset < 0 {
		return "negative offset"
	}

	return ""
}

type PostImg struct {
	Image  string
	PostID string
}
