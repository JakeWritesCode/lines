package stores

import (
	"github.com/stretchr/testify/assert"
	"lines/lines/store"
	"testing"
)

func TestUserPostgresStore_CreateUser(t *testing.T) {
	pgStore := NewUserPostgresStore()
	stores := []store.IntegrationTestStore{pgStore}
	store.IsolatedIntegrationTest(t, stores, func(t *testing.T) {
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

func TestUserPostgresStore_CreateUser_ValidationErrors(t *testing.T) {
	pgStore := NewUserPostgresStore()
	stores := []store.IntegrationTestStore{pgStore}
	store.IsolatedIntegrationTest(t, stores, func(t *testing.T) {
		user := User{
			Email:    "someemail.com",
			Password: "password",
		}
		validationErrors, err := pgStore.CreateUser(&user)
		assert.Nil(t, err)
		assert.NotEmpty(t, validationErrors)
	})
}

func TestUserPostgresStore_GetUserByEmail(t *testing.T) {
	pgStore := NewUserPostgresStore()
	stores := []store.IntegrationTestStore{pgStore}
	store.IsolatedIntegrationTest(t, stores, func(t *testing.T) {
		user := User{
			Name:     "Test User",
			Email:    "some@email.com",
			Password: "password",
		}
		validationErrors, err := pgStore.CreateUser(&user)
		assert.Nil(t, err)
		assert.Empty(t, validationErrors)
		assert.NotEqual(t, uint(0), user.ID)
		dbUser, err := pgStore.GetUserByEmail(user.Email)
		assert.Nil(t, err)
		assert.Equal(t, user.ID, dbUser.ID)
		assert.Equal(t, user.Name, dbUser.Name)
		assert.Equal(t, user.Email, dbUser.Email)
		assert.Equal(t, user.Password, dbUser.Password)
	})
}

func TestUserPostgresStore_GetUserByEmail_Error(t *testing.T) {
	pgStore := NewUserPostgresStore()
	stores := []store.IntegrationTestStore{pgStore}
	store.IsolatedIntegrationTest(t, stores, func(t *testing.T) {
		user, err := pgStore.GetUserByEmail("alala")
		assert.NotNil(t, err)
		assert.Nil(t, user)
	})
}

func TestUserPostgresStore_GetUserByID(t *testing.T) {
	pgStore := NewUserPostgresStore()
	stores := []store.IntegrationTestStore{pgStore}
	store.IsolatedIntegrationTest(t, stores, func(t *testing.T) {
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

func TestUserPostgresStore_GetUserByID_Error(t *testing.T) {
	pgStore := NewUserPostgresStore()
	stores := []store.IntegrationTestStore{pgStore}
	store.IsolatedIntegrationTest(t, stores, func(t *testing.T) {
		user, err := pgStore.GetUserByID(0)
		assert.NotNil(t, err)
		assert.Nil(t, user)
	})
}

func TestUserPostgresStore_UpdateUser(t *testing.T) {
	pgStore := NewUserPostgresStore()
	stores := []store.IntegrationTestStore{pgStore}
	store.IsolatedIntegrationTest(t, stores, func(t *testing.T) {
		user := User{
			Name:     "Test User",
			Email:    "some@email.com",
			Password: "password",
		}
		validationErrors, err := pgStore.CreateUser(&user)
		assert.Nil(t, err)
		assert.Empty(t, validationErrors)
		assert.NotEqual(t, uint(0), user.ID)
		user.Name = "Updated Name"
		validationErrors, err = pgStore.UpdateUser(&user)
		assert.Nil(t, err)
		assert.Empty(t, validationErrors)
		dbUser, err := pgStore.GetUserByID(user.ID)
		assert.Nil(t, err)
		assert.Equal(t, user.ID, dbUser.ID)
		assert.Equal(t, user.Name, dbUser.Name)
		assert.Equal(t, user.Email, dbUser.Email)
		assert.Equal(t, user.Password, dbUser.Password)
	})
}

func TestUserPostgresStore_UpdateUser_ValidationErrors(t *testing.T) {
	pgStore := NewUserPostgresStore()
	stores := []store.IntegrationTestStore{pgStore}
	store.IsolatedIntegrationTest(t, stores, func(t *testing.T) {
		user := User{
			Name:     "Test User",
			Email:    "alala",
			Password: "password",
		}
		validationErrors, err := pgStore.CreateUser(&user)
		assert.Nil(t, err)
		assert.Empty(t, validationErrors)
		assert.NotEqual(t, uint(0), user.ID)
		user.Email = ""
		validationErrors, err = pgStore.UpdateUser(&user)
		assert.NotEmpty(t, validationErrors)
		dbUser, err := pgStore.GetUserByID(user.ID)
		assert.Nil(t, err)
		assert.Equal(t, "Test User", dbUser.Name)
		assert.Equal(t, "alala", dbUser.Email)
		assert.Equal(t, "password", dbUser.Password)
	})
}

func TestUserPostgresStore_DeleteUser(t *testing.T) {
	pgStore := NewUserPostgresStore()
	stores := []store.IntegrationTestStore{pgStore}
	store.IsolatedIntegrationTest(t, stores, func(t *testing.T) {
		user := User{
			Name:     "Test User",
			Email:    "some@user.com",
			Password: "password",
		}
		validationErrors, err := pgStore.CreateUser(&user)
		assert.Nil(t, err)
		assert.Empty(t, validationErrors)
		assert.NotEqual(t, uint(0), user.ID)
		err = pgStore.DeleteUser(&user)
		assert.Nil(t, err)
		dbUser, err := pgStore.GetUserByID(user.ID)
		assert.NotNil(t, err)
		assert.Nil(t, dbUser)
	})
}
