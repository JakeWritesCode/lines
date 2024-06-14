package stores

import (
	"github.com/stretchr/testify/assert"
	"lines/lines/store"
	"testing"
)

func TestUserPostgresStore_CreateUser(t *testing.T) {

	pgStore := NewUserPostgresStore()
	store.IsolatedIntegrationTest(t, []store.IntegrationTestStore{pgStore}, func(t *testing.T) {
		user := User{
			Name:     "Test User",
			Email:    "some@email.com",
			Password: "password",
		}
		validationErrors, err := pgStore.CreateUser(&user)
		assert.Nil(t, err)
		assert.Empty(t, validationErrors)
		assert.NotEqual(t, uint(0), user.ID)
		dbUser, err := pgStore.GetUserByID(user.ID)
		assert.Nil(t, err)
		assert.Equal(t, user.ID, dbUser.ID)
		assert.Equal(t, user.Name, dbUser.Name)
		assert.Equal(t, user.Email, dbUser.Email)
		assert.Equal(t, user.Password, dbUser.Password)
	})
}
