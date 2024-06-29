package http

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserLogin_Validate_NoEmail(t *testing.T) {
	u := UserLogin{
		Password: "password",
	}
	err := u.Validate()
	assert.Equal(t, "Email is required.", err.Message[0])
}

func TestUserLogin_Validate_NoPassword(t *testing.T) {
	u := UserLogin{
		Email: "email",
	}
	err := u.Validate()
	assert.Equal(t, "Password is required.", err.Message[0])
}

func TestUserSignUp_Validate(t *testing.T) {
	u := UserSignUp{
		Name:     "name",
		Email:    "email",
		Password: "password",
	}
	err := u.Validate()
	assert.Equal(t, 0, len(err.Message))
}

func TestUserSignUp_Validate_NoEmail(t *testing.T) {
	u := UserSignUp{
		Name:     "name",
		Password: "password",
	}
	err := u.Validate()
	assert.Equal(t, "Email is required.", err.Message[0])
}

func TestUserSignUp_Validate_NoPassword(t *testing.T) {
	u := UserSignUp{
		Name:  "name",
		Email: "email",
	}
	err := u.Validate()
	assert.Equal(t, "Password is required.", err.Message[0])
}

func TestUserSignUp_Validate_NoName(t *testing.T) {
	u := UserSignUp{
		Email:    "email",
		Password: "password",
	}
	err := u.Validate()
	assert.Equal(t, "Name is required.", err.Message[0])
}
