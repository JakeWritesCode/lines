package stores

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestUser_Validate_NameIsEmpty(t *testing.T) {
	user := User{
		Name:     "",
		Email:    "some@email.com",
		Password: "password",
	}
	errors := user.Validate()
	assert.Equal(t, 1, len(errors))
	assert.Equal(t, "Name", errors[0].Field)
	assert.Equal(t, "Name is required", errors[0].Message)
}

func TestUser_Validate_EmailIsEmpty(t *testing.T) {
	user := User{
		Name:     "Name",
		Email:    "",
		Password: "password",
	}
	errors := user.Validate()
	assert.Equal(t, 1, len(errors))
	assert.Equal(t, "Email", errors[0].Field)
	assert.Equal(t, "Email is required", errors[0].Message)
}

func TestUser_Validate_PasswordIsEmpty(t *testing.T) {
	user := User{
		Name:     "Name",
		Email:    "some@email.com",
		Password: "",
	}
	errors := user.Validate()
	assert.Equal(t, 1, len(errors))
	assert.Equal(t, "Password", errors[0].Field)
	assert.Equal(t, "Password is required", errors[0].Message)
}
