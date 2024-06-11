package store

import "gorm.io/gorm"

// Store is a struct that contains a pointer to a gorm.DB instance
type Store struct {
	DB *gorm.DB
}

// NewStore is a function that returns a new Store instance
func NewStore(db *gorm.DB) *Store {
	return &Store{
		DB: db,
	}
}
