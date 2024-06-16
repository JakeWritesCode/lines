package domain

import (
	"github.com/stretchr/testify/assert"
	"lines/user/stores"
	"testing"
)

type MockUserStoreUserExists struct {
	stores.UserStoreInterface
}

func (m MockUserStoreUserExists) GetUserByEmail(email string) (*stores.User, error) {
	return &stores.User{
		Name: "Test User",
	}, nil
}

func TestUserForCreate_Validate_EmptyName(t *testing.T) {
	user := UserForCreate{
		Name:     "",
		Email:    "test@test.com",
		Password: "password",
	}
	validation, err := user.Validate(MockUserStoreUserExists{})
	assert.Nil(t, err)
	assert.NotEmpty(t, validation)
	assert.Equal(t, "name", validation[0].Field)
	assert.Equal(t, "name is required", validation[0].Errors[0])
}

func TestUserForCreate_Validate_EmptyEmail(t *testing.T) {
	user := UserForCreate{
		Name:     "Test User",
		Email:    "",
		Password: "password",
	}
	validation, err := user.Validate(MockUserStoreUserExists{})
	assert.Nil(t, err)
	assert.NotEmpty(t, validation)
	assert.Equal(t, "email", validation[0].Field)
	assert.Equal(t, "email is required", validation[0].Errors[0])
}

func TestUserForCreate_Validate_InvalidEmail(t *testing.T) {
	user := UserForCreate{
		Name:     "Test User",
		Email:    "sosig",
		Password: "password",
	}
	validation, err := user.Validate(MockUserStoreUserExists{})
	assert.Nil(t, err)
	assert.NotEmpty(t, validation)
	assert.Equal(t, "email", validation[0].Field)
	assert.Equal(t, "Invalid email", validation[0].Errors[0])
}

func TestUserForCreate_Validate_EmptyPassword(t *testing.T) {
	user := UserForCreate{
		Name:     "Test User",
		Email:    "some@email.com",
		Password: "",
	}
	validation, err := user.Validate(MockUserStoreUserExists{})
	assert.Nil(t, err)
	assert.NotEmpty(t, validation)
	assert.Equal(t, "password", validation[0].Field)
	assert.Equal(t, "password is required", validation[0].Errors[0])
}

func TestUserForCreate_Validate_EmailExists(t *testing.T) {
	user := UserForCreate{
		Name:     "Test User",
		Email:    "some@email.com",
		Password: "password",
	}
	validation, err := user.Validate(MockUserStoreUserExists{})
	assert.Nil(t, err)
	assert.NotEmpty(t, validation)
	assert.Equal(t, "email", validation[0].Field)
	assert.Equal(t, "Email is already in use.", validation[0].Errors[0])
}

type MockUserStoreUserDoesNotExist struct {
	stores.UserStoreInterface
}

func (m MockUserStoreUserDoesNotExist) GetUserByEmail(email string) (*stores.User, error) {
	return nil, nil
}

func TestUserForCreate_Validate_Valid(t *testing.T) {
	user := UserForCreate{
		Name:     "Test User",
		Email:    "some@email.com",
		Password: "password",
	}
	validation, err := user.Validate(MockUserStoreUserDoesNotExist{})
	assert.Nil(t, err)
	assert.Empty(t, validation)
}

type MockUserStoreError struct {
	stores.UserStoreInterface
}

func (m MockUserStoreError) GetUserByEmail(email string) (*stores.User, error) {
	return nil, assert.AnError
}

func TestUserForCreate_Validate_Error(t *testing.T) {
	user := UserForCreate{
		Name:     "Test User",
		Email:    "some@user.com",
		Password: "password",
	}
	validation, err := user.Validate(MockUserStoreError{})
	assert.NotNil(t, err)
	assert.Empty(t, validation)
}
