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

// PostgresDB is a PostgresDB.
type PostgresDB struct {
	config PostgresDBConfig
	DB     *gorm.DB
}

// PostgresStore is a struct that contains an initialized PostgresDB instance.
type PostgresStore struct {
	Postgres *PostgresDB
}
