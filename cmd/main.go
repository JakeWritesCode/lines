package main

import (
	"fmt"
	"lines/internal"
	"lines/lines/app"
	"lines/lines/http"
)

// Main is the entry point for the application.
func main() {
	apps := []app.App{
		// Add your apps here.
	}

	appConfig := internal.NewConfig()
	httpEngine := http.CreateEngine(*appConfig)
	// TODO: Here we're going to initialise sentry, datadog and other app based stuff.

	// Next, we initialise each of our apps. Each app then initialises its own dependencies.
	for _, a := range apps {
		err := a.Initialise(appConfig)
		if err != nil {
			appConfig.Logger.Error(
				"main",
				"main",
				fmt.Sprintf("Failed to initialise app: %s", err.Error()),
			)
			panic(err)
		}
	}

}
