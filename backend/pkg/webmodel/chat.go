package webmodel

import (
	"time"
)

type ChatMessage struct {
	MessageID string `json:"messageID,omitempty"`
	UserID    string `json:"userID,omitempty"` // server uses current user to identify the message's author and chat which the user is writing to
	// TODO del GroupID    string    `json:"groupID,omitempty"`
	UserName   string    `json:"userName,omitempty"`
	Content    string    `json:"content"`
	DateCreate time.Time `json:"dateCreate,omitempty"`
	Images     []string  `json:"images,omitempty"`
	GenericID  string    `json:"genericID,omitempty"` // generic ID for frontend, either recipient userID for private chat or groupID for group chat
}

func (m *ChatMessage) Validate() string {
	if m.DateCreate.Before(time.Date(2023, time.September, 1, 0, 0, 0, 0, time.UTC)) {
		return "Date is too old"
	}
	if IsEmpty(m.Content) {
		return "text is missing"
	}

	return ""
}
