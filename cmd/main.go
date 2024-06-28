package main

import (
	"fmt"
	"lines/internal"
	"lines/lines/app"
	"lines/lines/http"
	"lines/user"
	"strconv"
)

// Main is the entry point for the application.
func main() {
	userApp := user.NewUserApp()
	apps := []app.App{
		&userApp,
	}
	config := internal.NewConfig()
	httpEngine := http.CreateEngine(config)

	MainHandler(apps, config, httpEngine)
}

// MainHandler is the main handler for the application.
func MainHandler(
	apps []app.App,
	config *internal.MainConfig,
	httpEngine http.HttpEngine,
) {
	// TODO: Here we're going to initialise sentry, datadog and other app based stuff.

	// Next, we initialise each of our apps. Each app then initialises its own dependencies.
	for _, a := range apps {
		err := a.Initialise()
		if err != nil {
			config.Logger.Fatal(
				"main",
				"main",
				fmt.Sprintf("Failed to initialise app: %s", err.Error()),
			)
		}
		a.RegisterHTTPRoutes(httpEngine)
	}

	// Start the server.
	err := httpEngine.Run(":" + strconv.Itoa(config.HTTPPort))
	if err != nil {
		config.Logger.Fatal(
			"main",
			"main",
			fmt.Sprintf("Failed to start server: %s", err.Error()),
		)
	}
}
