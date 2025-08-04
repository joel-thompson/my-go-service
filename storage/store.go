package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
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

// ListItems retrieves a paginated list of items from the database
func (s *Store) ListItems(ctx context.Context, req ListItemsRequest) (*ListItemsResponse, error) {
	// Set default values for pagination
	if req.Limit <= 0 {
		req.Limit = 10 // Default to 10 items per page
	}
	if req.Limit > 100 {
		req.Limit = 100 // Maximum 100 items per page
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	// Get total count
	var total int
	err := s.db.GetContext(ctx, &total, countItemsQuery)
	if err != nil {
		return nil, err
	}

	// Get items
	var items []Item
	err = s.db.SelectContext(ctx, &items, listItemsQuery, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	return &ListItemsResponse{
		Items:  items,
		Total:  total,
		Limit:  req.Limit,
		Offset: req.Offset,
	}, nil
}

// GetItem retrieves a single item by ID
func (s *Store) GetItem(ctx context.Context, id uuid.UUID) (*Item, error) {
	var item Item
	err := s.db.GetContext(ctx, &item, getItemQuery, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("item not found")
		}
		return nil, err
	}
	return &item, nil
}

// UpdateItem updates an existing item
func (s *Store) UpdateItem(ctx context.Context, id uuid.UUID, req UpdateItemRequest) (*Item, error) {
	var item Item
	err := s.db.GetContext(ctx, &item, updateItemQuery, id, req.Name, req.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("item not found")
		}
		return nil, err
	}
	return &item, nil
}

// DeleteItem deletes an item by ID
func (s *Store) DeleteItem(ctx context.Context, id uuid.UUID) (*Item, error) {
	var item Item
	err := s.db.GetContext(ctx, &item, deleteItemQuery, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("item not found")
		}
		return nil, err
	}
	return &item, nil
}
