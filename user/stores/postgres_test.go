package stores

import (
	"github.com/stretchr/testify/assert"
	"lines/lines/store"
	"testing"
)

type MockModel struct {
}

func (m MockModel) Validate() []store.ModelValidationError {
	return nil
}

func TestNewUserPostgresStore_Integration(t *testing.T) {
	pgStore := NewUserPostgresStore()
	assert.Equal(t, pgStore.PostgresStore.Config.AppName, "USER")
	assert.NotNil(t, pgStore.Postgres)
	assert.NotNil(t, pgStore.Logger)
}
