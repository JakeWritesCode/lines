package store

import "gorm.io/gorm"

// User is a struct that contains the fields of a user
type User struct {
	gorm.Model
	Email     string `gorm:"uniqueIndex"`
	Password  string
	FirstName string
	LastName  string
}

// UserObjectPermission is a struct that contains the fields of a user object permission.
type UserObjectPermission struct {
	gorm.Model
	UserID     uint
	ObjectID   uint
	ObjectType string
	Permission string
}
