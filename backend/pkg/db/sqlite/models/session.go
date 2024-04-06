package models

import "time"

type Session struct {
	ID      string    `json:"id"`
	Expired time.Time `json:"expired,omitempty"`
}
