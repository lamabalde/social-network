package models

import (
	"fmt"
	"time"
)

type Chat struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Type     int            `json:"type"` // 0 - private,  1 - group
	Messages []*ChatMessage `json:"messages,omitempty"`
}

type ChatMessage struct {
	ID         string    `json:"id"`
	UserID     string    `json:"userID,omitempty"`
	UserName   string    `json:"userName,omitempty"`
	Content    string    `json:"content"`
	DateCreate time.Time `json:"dateCreate,omitempty"`
	Images     []string  `json:"-"`
}

func (m ChatMessage) String() string {
	return fmt.Sprintf("ID:%s\n    Author:    %s (id: %s)\n    Content: %s\n    DataCreate: %s \n",
		m.ID, m.UserName, m.UserID, m.Content, m.DateCreate.String())
}

func (c *Chat) String() string {
	if c == nil {
		return "nil"
	}

	messages := ""
	for _, m := range c.Messages {
		messages += fmt.Sprintf("%s\n--------\n", m)
	}
	return fmt.Sprintf("chat id: %s | name: %s | type: %d  \n   Messages:\n%s\n\n", c.ID, c.Name, c.Type, messages)
}
