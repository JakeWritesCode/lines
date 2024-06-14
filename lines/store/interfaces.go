package store

import (
	"gorm.io/gorm"
	"lines/lines/logging"
)

// PostgresDBConfig is the configuration for a PostgresDB.
type PostgresDBConfig struct {
	Logger           logging.Logger
	ConnectionString string
	AppName          string
}

// PostgresStore is a struct that contains an initialized PostgresDB instance.
type PostgresStore struct {
	Postgres *gorm.DB
	Config   PostgresDBConfig
}

type ModelValidationError struct {
	Field   string
	Message string
}

type PostgresModel interface {
	Validate() []ModelValidationError
}
