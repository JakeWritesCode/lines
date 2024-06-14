package stores

import (
	"gorm.io/gorm"
	"lines/lines/store"
)

type User struct {
	gorm.Model
	Name     string
	Email    string
	Password string
}

func (u User) Validate() []store.ModelValidationError {
	var errors []store.ModelValidationError
	// TODO: Find a way to create generic validation functions
	if u.Name == "" {
		errors = append(errors, store.ModelValidationError{Field: "Name", Message: "Name is required"})
	}
	if u.Email == "" {
		errors = append(errors, store.ModelValidationError{Field: "Email", Message: "Email is required"})
	}
	if u.Password == "" {
		errors = append(errors, store.ModelValidationError{Field: "Password", Message: "Password is required"})
	}
	return errors
}
