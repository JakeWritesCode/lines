package stores

import (
	"lines/lines/logging"
	"lines/lines/store"
)

// TUserPostgresStore is an interface for a UserPostgresStore.
type UserPostgresStoreInterface interface {
	CreateUser(user *User) ([]store.ModelValidationError, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id uint) (*User, error)
	UpdateUser(user *User) ([]store.ModelValidationError, error)
}

// UserPostgresStore is a struct that contains an initialized PostgresStore instance.
type UserPostgresStore struct {
	*store.PostgresStore
	Logger logging.Logger
}

// NewUserPostgresStore is a function that returns a new UserPostgresStore instance.
func NewUserPostgresStore(models []store.PostgresModel) *UserPostgresStore {
	appName := "USER"
	config := store.CreatePostgresDBConfig(appName)
	db := store.CreatePostgresDB(*config, models)
	return &UserPostgresStore{
		PostgresStore: db,
		Logger:        config.Logger,
	}
}
