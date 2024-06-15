package store

import (
	"database/sql"
	"errors"
	"gorm.io/gorm"
	"lines/lines/logging"
)

// PostgresDBConfig is the configuration for a PostgresDB.
type PostgresDBConfig struct {
	TestRunner       bool
	Logger           logging.Logger
	ConnectionString string
	AppName          string
}

type PostgresStoreInterface interface {
	BeginTransaction() error
	RollbackTransaction() error
	Models() []PostgresModel
}

// PostgresStore is a struct that contains an initialized PostgresDB instance.
type PostgresStore struct {
	Postgres GormInstanceInterface
	Config   PostgresDBConfig
}

func (s *PostgresStore) BeginTransaction() error {
	if !s.Config.TestRunner {
		return errors.New("cannot start transaction in non-test environment")
	}
	s.Postgres = s.Postgres.Begin()
	return nil
}

func (s *PostgresStore) RollbackTransaction() error {
	s.Postgres = s.Postgres.Rollback()
	return nil
}

func (s *PostgresStore) Models() []PostgresModel {
	return []PostgresModel{}
}

type ModelValidationError struct {
	Field   string
	Message string
}

type PostgresModel interface {
	Validate() []ModelValidationError
}

type IntegrationTestStore interface {
	BeginTransaction() error
	RollbackTransaction() error
}

type GormInstanceInterface interface {
	Create(value interface{}) (tx *gorm.DB)
	Begin(opts ...*sql.TxOptions) *gorm.DB
	Rollback() *gorm.DB
	AutoMigrate(dst ...interface{}) error
	Where(query interface{}, args ...interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	Delete(value interface{}, conds ...interface{}) *gorm.DB
}
