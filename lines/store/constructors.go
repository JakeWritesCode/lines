package store

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"lines/lines/logging"
	"lines/lines/utils"
)

// CreatePostgresDBConfig creates a new PostgresDBConfig instance.
func CreatePostgresDBConfig(AppName string) *PostgresDBConfig {
	logLevel := utils.GetEnvOrDefault("LOG_LEVEL", "info", "string").(string)
	testRunner := utils.GetEnvOrDefault("TEST_RUNNER", "false", "bool").(bool)
	logger := logging.NewLogrusHandler(logLevel)
	dbName := ""
	if testRunner {
		dbName = fmt.Sprintf("%v_POSTGRES_URL_TEST", AppName)
	} else {
		dbName = fmt.Sprintf("%v_POSTGRES_URL", AppName)
	}
	connString := utils.GetEnvOrDefault(
		dbName,
		"NODEFAULT",
		"string",
	).(string)
	return &PostgresDBConfig{
		Logger:           logger,
		ConnectionString: connString,
		AppName:          AppName,
		TestRunner:       testRunner,
	}
}

// CreatePostgresDB creates a new PostgresDB instance.
// It connects to the database using the provided configuration.
// It also migrates the provided models to the database.
func CreatePostgresDB(config PostgresDBConfig, models []PostgresModel) *gorm.DB {
	db, err := gorm.Open(postgres.Open(config.ConnectionString), &gorm.Config{})
	if err != nil {
		config.Logger.Fatal(
			config.AppName,
			"CreatePostgresDB",
			fmt.Sprintf("Failed to connect to the database: %v", err),
		)
		return nil
	}
	for _, model := range models {
		err = db.AutoMigrate(model)
		if err != nil {
			config.Logger.Fatal(
				config.AppName,
				"CreatePostgresDB",
				fmt.Sprintf("Failed to migrate the database: %v", err),
			)
		}
	}
	return db
}
