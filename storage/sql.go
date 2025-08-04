package storage

import (
	"github.com/google/uuid"
	"time"
)

// Item represents an item in the database
type Item struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description *string   `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// CreateItemRequest represents the request payload for creating an item
type CreateItemRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

// ListItemsRequest represents pagination parameters for listing items
type ListItemsRequest struct {
	Limit  int `form:"limit" json:"limit"`
	Offset int `form:"offset" json:"offset"`
}

// ListItemsResponse represents the response for listing items
type ListItemsResponse struct {
	Items  []Item `json:"items"`
	Total  int    `json:"total"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

const (
	createItemQuery = `
		INSERT INTO items (name, description)
		VALUES ($1, $2)
		RETURNING id, name, description, created_at, updated_at
	`

	listItemsQuery = `
		SELECT id, name, description, created_at, updated_at
		FROM items
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	countItemsQuery = `
		SELECT COUNT(*)
		FROM items
	`
)
