package webmodel

import (
	"time"
)

type Group struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DateCreate  time.Time `json:"dateCreate,omitempty"`
}

func (g Group) Validate() string {
	// if g.DateCreate.Before(time.Date(2023, time.September, 1, 0, 0, 0, 0, time.UTC)) {
	// 	return "Date is too old"
	// }

	if IsEmpty(g.Title) {
		return "Title is missing"
	}

	if IsEmpty(g.Description) {
		return "Description is missing"
	}

	return ""
}

type GroupUser struct {
	GroupID  string `json:"groupID"`
	User     string `json:"user"` // ID or search string
}

func (ug GroupUser) Validate() string {
	if IsEmpty(ug.GroupID) {
		return "GroupID is missing"
	}

	if IsEmpty(ug.User) {
		return "User is missing"
	}

	return ""
}

type Event struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	GroupID     string    `json:"groupID"`
	DateCreate  time.Time `json:"dateCreate,omitempty"`
	DateEvent   time.Time `json:"dateEvent"`
}

func (ug Event) Validate() string {
	if IsEmpty(ug.Title) {
		return "Title is missing"
	}

	if IsEmpty(ug.Description) {
		return "Description is missing"
	}

	if IsEmpty(ug.GroupID) {
		return "GroupID is missing"
	}

	if ug.DateEvent.Before(time.Date(202, time.January, 1, 0, 0, 0, 0, time.UTC)) {
		return "DateEvent is too early"
	}
	return ""
}

const NUM_OF_EVENT_OPTIONS = 2
const (
	EVENT_NOT_GOING = iota
	EVENT_GOING
)

type UserOptionForEvent struct {
	EventID string `json:"eventID"`
	Option  int    `json:"option"`
}

func (uo UserOptionForEvent) Validate() string {
	if IsEmpty(uo.EventID) {
		return "EventID is missing"
	}

	if uo.Option < 0 || uo.Option >= NUM_OF_EVENT_OPTIONS {
		return "Wrong option"
	}

	return ""
}
