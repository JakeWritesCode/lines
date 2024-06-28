package store

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCreatePostgresDBConfig(t *testing.T) {
	config := CreatePostgresDBConfig("Testapp")
	assert.NotNil(t, config)
	assert.NotNil(t, config.Logger)
	assert.NotEmpty(t, config.ConnectionString)
	assert.Contains(t, config.ConnectionString, "NODEFAULT")
	assert.Equal(t, "Testapp", config.AppName)
	assert.True(t, config.TestRunner)
}

func TestCreatePostgresDBConfig_nonDefaults(t *testing.T) {
	// Set the environment variables
	assert.Nil(t, os.Setenv("LOG_LEVEL", "debug"))
	assert.Nil(t, os.Setenv("TEST_RUNNER", "true"))
	assert.Nil(t, os.Setenv("Testapp_POSTGRES_URL_TEST", "postgres://localhost:5432/testapp"))
	config := CreatePostgresDBConfig("Testapp")
	assert.NotNil(t, config)
	assert.NotNil(t, config.Logger)
	assert.NotEmpty(t, config.ConnectionString)
	assert.Equal(t, "postgres://localhost:5432/testapp", config.ConnectionString)
	assert.Equal(t, "Testapp", config.AppName)
	os.Clearenv()
}
