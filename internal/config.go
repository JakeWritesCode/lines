package internal

// MainConfig is the main configuration struct for the whole monolith.
type MainConfig struct {
	LocalDev    bool
	TestRunner  bool
	UseSSL      bool
	SiteDomain  string
	LogLevel    string
	Logger      Logger
	CORSOrigins []string
	SentryDSN   string
	HTTPPort    int
}

// NewConfig creates a new MainConfig struct, reading from environment variables.
func NewConfig() *MainConfig {
	localDev := GetEnvOrDefault("LOCAL_DEV", "true", "bool").(bool)
	testRunner := GetEnvOrDefault("TEST_RUNNER", "false", "bool").(bool)
	useSSL := GetEnvOrDefault("USE_SSL", "false", "bool").(bool)
	config := &MainConfig{
		LocalDev:    localDev,
		TestRunner:  testRunner,
		UseSSL:      useSSL && !localDev,
		SiteDomain:  GetEnvOrDefault("SITE_DOMAIN", "localhost", "string").(string),
		LogLevel:    GetEnvOrDefault("LOG_LEVEL", "info", "string").(string),
		CORSOrigins: GetEnvOrDefault("CORS_ORIGINS", "http://localhost", "[]string").([]string),
		SentryDSN:   GetEnvOrDefault("SENTRY_DSN", "", "string").(string),
		HTTPPort:    GetEnvOrDefault("HTTP_PORT", "8080", "int").(int),
	}
	config.Logger = NewLogrusHandler(config.LogLevel)
	return config
}
