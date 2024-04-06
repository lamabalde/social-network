package models

type SearchResult struct {
	EntityType string `json:"entityType"` // groups or profile
	ID         string `json:"id"`
	Name       string `json:"name"`
}
