package store

import (
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

type MockGormInstance struct {
	GormInstanceInterface
}

func (m *MockGormInstance) Begin(...*sql.TxOptions) *gorm.DB {
	return &gorm.DB{}
}

func TestPostgresStore_BeginTransaction(t *testing.T) {
	mockGormInstance := &MockGormInstance{}
	store := &PostgresStore{
		Postgres: mockGormInstance,
	}

	err := store.BeginTransaction()
	assert.Nil(t, err)
}

func (m *MockGormInstance) Rollback() *gorm.DB {
	return &gorm.DB{}
}

func TestPostgresStore_RollbackTransaction(t *testing.T) {
	mockGormInstance := &MockGormInstance{}
	store := &PostgresStore{
		Postgres: mockGormInstance,
	}

	err := store.RollbackTransaction()
	assert.Nil(t, err)
}

func TestPostgresStore_Models(t *testing.T) {
	store := &PostgresStore{}

	models := store.Models()
	assert.Len(t, models, 0)
}

func TestPostgresStore_RecordNotFound(t *testing.T) {
	store := &PostgresStore{}

	isError := store.RecordNotFound(gorm.ErrRecordNotFound)
	assert.True(t, isError)
	isError = store.RecordNotFound(errors.New("some error"))
	assert.False(t, isError)
}
