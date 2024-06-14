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
	Postgres store.GormInstanceInterface
	Logger   logging.Logger
	Config   store.PostgresDBConfig
}

func (s *UserPostgresStore) Models() []store.PostgresModel {
	return []store.PostgresModel{
		User{},
	}
}

// NewUserPostgresStore is a function that returns a new UserPostgresStore instance.
func NewUserPostgresStore() *UserPostgresStore {
	appName := "USER"
	config := store.CreatePostgresDBConfig(appName)
	postgresStore := &UserPostgresStore{}
	db := store.CreatePostgresDB(*config, postgresStore.Models())
	return &UserPostgresStore{
		Postgres: db,
		Logger:   config.Logger,
		Config:   *config,
	}
}
