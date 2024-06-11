package internal

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewConfig_SetsBasicValues(t *testing.T) {
	envMap := map[string]string{
		"LOCAL_DEV":    "false",
		"TEST_RUNNER":  "false",
		"USE_SSL":      "true",
		"SITE_DOMAIN":  "test.com",
		"CORS_ORIGINS": "http://test.com,http://test2.com",
		"SENTRY_DSN":   "sentry",
		"LOG_LEVEL":    "info",
	}
	for k, v := range envMap {
		err := os.Setenv(k, v)
		if err != nil {
			t.Error("Failed to set environment variable")
		}
	}
	config := NewConfig()
	assert.False(t, config.LocalDev)
	assert.False(t, config.TestRunner)
	assert.True(t, config.UseSSL)
	assert.Equal(t, "test.com", config.SiteDomain)
	assert.Equal(t, "sentry", config.SentryDSN)
	assert.Equal(t, []string{"http://test.com", "http://test2.com"}, config.CORSOrigins)
	assert.Equal(t, "info", config.LogLevel)
}

func TestNewConfig_SetsHttpWhenOnLocalDev(t *testing.T) {
	envMap := map[string]string{
		"LOCAL_DEV": "true",
		"USE_SSL":   "true",
	}
	for k, v := range envMap {
		err := os.Setenv(k, v)
		if err != nil {
			t.Error("Failed to set environment variable")
		}
	}
	config := NewConfig()
	assert.True(t, config.LocalDev)
	assert.False(t, config.UseSSL)
}

func TestNewConfig_InitialisesLogger(t *testing.T) {
	envMap := map[string]string{
		"LOG_LEVEL": "info",
	}
	for k, v := range envMap {
		err := os.Setenv(k, v)
		if err != nil {
			t.Error("Failed to set environment variable")
		}
	}
	config := NewConfig()
	assert.NotNil(t, config.Logger)
}
