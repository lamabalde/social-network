package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type Post struct {
	ID               string     `json:"id,omitempty"`
	Theme            string     `json:"theme"`
	Content          Content    `json:"content"`
	Categories       []string   `json:"categories"`
	Comments         []*Comment `json:"comments,omitempty"`
	CommentsQuantity int        `json:"commentsQuantity,omitempty"`
	GroupID          string     `json:"groupID,omitempty"`
	Image            string     `json:"image,omitempty"`
	Privacy          int        `json:"privacy,omitempty"` // 0 is public, 1 is private, 2 is friends
}

func (p *Post) String() string {
	if p == nil {
		return "nil"
	}
	return fmt.Sprintf("id: %s | Theme: %s\nContent: \n%s\nCategories: \n%v\nComments(%d): \n%v\nGroupID: %s\n",
		p.ID, p.Theme, p.Content.String(), p.Categories, p.CommentsQuantity, p.Comments, p.GroupID)
}

type Comment struct {
	ID      string  `json:"id,omitempty"`
	PostID  string  `json:"postID"`
	Content Content `json:"content"`
}

func (c *Comment) String() string {
	if c == nil {
		return "nil"
	}
	return fmt.Sprintf("comment id: %s | Comment Message: \n%s\n", c.ID, c.Content.String())
}

func (c *Comment) UnmarshalJSON(data []byte) error {
	type tmpComment Comment
	var tmp tmpComment
	err := json.Unmarshal(data, &tmp)
	*c = Comment(tmp)
	return err
}

type Content struct {
	UserID       string        `json:"userID,omitempty"`
	UserName     string        `json:"userName,omitempty"`
	Text         string        `json:"text"`
	DateCreate   time.Time     `json:"dateCreate,omitempty"`
	Likes        []int         `json:"likes,omitempty"` // index 0 keeps number of dislikes, index 1 keeps number of likes
	Images       []string      `json:"-"`
	Image        string        `json:"image,omitempty"`
	UserReaction UserReactions `json:"userReaction,omitempty"` //-1 => no reaction
}

func (c *Content) String() string {
	if c == nil {
		return "nil"
	}
	return fmt.Sprintf("   Author: \n   ID: %v\n   username: %v\n   Text: %s\n   DataCreate: %s | Likes: %#v | UserReaction: %#v\n",
		c.UserID, c.UserName, c.Text, c.DateCreate.String(), c.Likes, c.UserReaction)
}

func (c *Content) UnmarshalJSON(data []byte) error {
	type tmpContent Content
	var tmp tmpContent
	err := json.Unmarshal(data, &tmp)
	*c = Content(tmp)
	return err
}
