package stores

import (
	"gorm.io/gorm"
	"lines/lines/logging"
	"lines/lines/store"
)

type UserPostgresStore struct {
	*store.PostgresStore
	Logger logging.Logger
}

func NewUserPostgresStore(models []gorm.Model) *UserPostgresStore {
	appName := "USER"
	config := store.CreatePostgresDBConfig(appName)
	db := store.CreatePostgresDB(*config, models)
	return &UserPostgresStore{
		PostgresStore: &store.PostgresStore{Postgres: db},
		Logger:        config.Logger,
	}
}
