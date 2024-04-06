package models

import (
	"fmt"
	"time"
)

type Group struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatorID   string     `json:"creator_id"`
	DateCreate  time.Time  `json:"dateCreate"`
	Members     []UserBase `json:"members"`
	Posts       []*Post    `json:"posts"`
}

func (g *Group) String() string {
	if g == nil {
		return "nil"
	}
	members := ""
	for _, m := range g.Members {
		members += fmt.Sprintf("- %s\n", m)
	}

	posts := ""
	for _, p := range g.Posts {
		posts += fmt.Sprintf("%s\n------------------\n", p)
	}

	return fmt.Sprintf("\nGroup: '%s' (id: '%s') -----------\n    Description:'%s'\n    Created:%s by user id:%s\n    Members:\n%s\n    Posts:\n%s",
		g.Name, g.ID, g.Description, g.DateCreate, g.CreatorID, members, posts)
}
