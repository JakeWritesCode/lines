package stores

import (
	"github.com/stretchr/testify/assert"
	store2 "lines/lines/store"
	"testing"
)

func TestNewUserStore(t *testing.T) {
	store := NewUserStore()
	assert.NotNil(t, store)
	assert.NotNil(t, store.UserPostgresStore)
	assert.NotNil(t, store.UserPostgresStore.PostgresStore)
	store2.IsolatedIntegrationTest(t, []store2.IntegrationTestStore{store.UserPostgresStore.PostgresStore}, func(t *testing.T) {
		user := User{
			Name:     "Test User",
			Email:    "Barry",
			Password: "password",
		}
		validationErrors, err := store.UserPostgresStore.CreateUser(&user)
		assert.Nil(t, err)
		assert.Empty(t, validationErrors)
		assert.NotEqual(t, uint(0), user.ID)
	})
}
