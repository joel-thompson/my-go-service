package storage

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// Store handles all database operations
type Store struct {
	db *sqlx.DB
}

// New creates a new Store instance
func New(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

// CreateItem creates a new item in the database
func (s *Store) CreateItem(ctx context.Context, req CreateItemRequest) (*Item, error) {
	var item Item
	err := s.db.GetContext(ctx, &item, createItemQuery, req.Name, req.Description)
	if err != nil {
		return nil, err
	}
	return &item, nil
}
