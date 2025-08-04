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

const (
	createItemQuery = `
		INSERT INTO items (name, description)
		VALUES ($1, $2)
		RETURNING id, name, description, created_at, updated_at
	`
)
