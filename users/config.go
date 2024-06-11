package internal

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type GlobalConfig struct {
	LocalDev                   bool
	TestRunner                 bool
	UseSSL                     bool
	SiteDomain                 string
	SecretKey                  []byte
	TokenExpirationTimeMinutes int
	Port                       int
	LogLevel                   string
	Logger                     Logger
	CORSOrigins                []string
	SentryDSN                  string
	PostgresURL                string
	PostgresPort               int
	PostgresUsername           string
	PostgresPassword           string
	PostgresDBName             string
	RedisURL                   string
	RedisPassword              string
}

func (c *Config) RedisConnStringSplitter() (string, string) {
	url := strings.Split(c.RedisURL, "@")[1]
	password := strings.Split(c.RedisURL, "@")[0]
	password = strings.Split(password, "redis://:")[1]
	return url, password
}

func (c *Config) GeneratePostgresConnString() string {
	return fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=disable",
		c.PostgresURL,
		c.PostgresUsername,
		c.PostgresPassword,
		c.PostgresDBName,
		c.PostgresPort,
	)
}

func NewConfig() *Config {
	localDev := GetEnvOrDefault("LOCAL_DEV", "true", "bool").(bool)
	testRunner := GetEnvOrDefault("TEST_RUNNER", "false", "bool").(bool)
	allowHTTP := GetEnvOrDefault("ALLOW_HTTP", "false", "bool").(bool)
	RedisConnString := GetEnvOrDefault("REDIS_CONN_STRING", "", "string").(string)
	config := &Config{
		LocalDev:                   localDev,
		TestRunner:                 testRunner,
		UseSSL:                     !allowHTTP && !localDev,
		SiteDomain:                 GetEnvOrDefault("SITE_DOMAIN", "localhost", "string").(string),
		SecretKey:                  []byte(GetEnvOrDefault("SECRET_KEY", "my_secret_key", "string").(string)),
		TokenExpirationTimeMinutes: GetEnvOrDefault("TOKEN_EXPIRATION_TIME_MINUTES", "60", "int").(int),
		Port:                       GetEnvOrDefault("PORT", "8080", "int").(int),
		LogLevel:                   GetEnvOrDefault("LOG_LEVEL", "info", "string").(string),
		CORSOrigins:                GetEnvOrDefault("CORS_ORIGINS", "", "[]string").([]string),
		SentryDSN:                  GetEnvOrDefault("SENTRY_DSN", "", "string").(string),
		PostgresURL:                GetEnvOrDefault("POSTGRES_URL", "localhost", "string").(string),
		PostgresPort:               GetEnvOrDefault("POSTGRES_PORT", "5432", "int").(int),
		PostgresUsername:           GetEnvOrDefault("POSTGRES_USERNAME", "postgres", "string").(string),
		PostgresPassword:           GetEnvOrDefault("POSTGRES_PASSWORD", "postgres", "string").(string),
		PostgresDBName:             GetEnvOrDefault("POSTGRES_DB_NAME", "garduino", "string").(string),
		RedisURL:                   GetEnvOrDefault("REDIS_URL", "localhost:6379", "string").(string),
		RedisPassword:              GetEnvOrDefault("REDIS_PASSWORD", "", "string").(string),
	}
	if RedisConnString != "" {
		config.RedisURL, config.RedisPassword = config.RedisConnStringSplitter()
	}
	return config
}

// GetEnvOrDefault returns the value of the environment variable named by the key.
// If the environment variable is empty, the defaultValue is returned.
// If the environment variable is not empty, the value is coerced to the expectedType.
// If the value cannot be coerced to the expectedType, panic.
func GetEnvOrDefault(key string, defaultValue string, expectedType string) interface{} {
	Value := os.Getenv(key)
	if Value == "" {
		Value = defaultValue
	}
	switch expectedType {
	case "bool":
		rVal, ok := coerceBool(Value)
		if !ok {
			panic("invalid bool value for " + key)
		}
		return rVal
	case "int":
		rVal, ok := coerceInt(Value)
		if !ok {
			panic("invalid int value for " + key)
		}
		return rVal
	case "string":
		return Value
	case "[]string":
		return coerceSliceOfStrings(Value)
	default:
		panic("invalid expectedType for " + key)
	}
}

// coerceBool converts a string to a bool.
// If the string is not a valid bool, the second return value is false.
func coerceBool(value string) (bool, bool) {
	if value == "true" {
		return true, true
	}
	if value == "false" {
		return false, true
	}
	return false, false
}

// coerceInt converts a string to an int.
// If the string is not a valid int, the second return value is false.
func coerceInt(value string) (int, bool) {
	convertedInt, err := strconv.Atoi(value)
	if err != nil {
		return 0, false
	}
	return convertedInt, true
}

func coerceSliceOfStrings(value string) []string {
	convertedSlice := []string{}
	if value == "" {
		return convertedSlice
	}
	return strings.Split(value, ",")
}
