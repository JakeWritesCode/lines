package internal

import (
	"lines/lines/logging"
	"lines/lines/utils"
)

// MainConfig is the main configuration struct for the whole monolith.
type MainConfig struct {
	LocalDev    bool
	TestRunner  bool
	UseSSL      bool
	SiteDomain  string
	LogLevel    string
	Logger      logging.Logger
	CORSOrigins []string
	SentryDSN   string
	HTTPPort    int
}

// NewConfig creates a new MainConfig struct, reading from environment variables.
func NewConfig() *MainConfig {
	localDev := utils.GetEnvOrDefault("LOCAL_DEV", "true", "bool").(bool)
	testRunner := utils.GetEnvOrDefault("TEST_RUNNER", "false", "bool").(bool)
	useSSL := utils.GetEnvOrDefault("USE_SSL", "false", "bool").(bool)
	config := &MainConfig{
		LocalDev:    localDev,
		TestRunner:  testRunner,
		UseSSL:      useSSL && !localDev,
		SiteDomain:  utils.GetEnvOrDefault("SITE_DOMAIN", "localhost", "string").(string),
		LogLevel:    utils.GetEnvOrDefault("LOG_LEVEL", "info", "string").(string),
		CORSOrigins: utils.GetEnvOrDefault("CORS_ORIGINS", "http://localhost", "[]string").([]string),
		SentryDSN:   utils.GetEnvOrDefault("SENTRY_DSN", "", "string").(string),
		HTTPPort:    utils.GetEnvOrDefault("HTTP_PORT", "8080", "int").(int),
	}
	config.Logger = logging.NewLogrusHandler(config.LogLevel)
	return config
}
