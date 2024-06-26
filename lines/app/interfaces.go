package app

import (
	linesHttp "lines/lines/http"
)

// App is the interface that all apps must implement.
type App interface {
	// Initialise is called to initialise the app.
	Initialise() error
	// RegisterHTTPRoutes is called to register the app's HTTP routes.
	RegisterHTTPRoutes(engine linesHttp.HttpEngine)
	// RegisterGRPCServices is called to register the app's gRPC services.
	RegisterGRPCServices() error
}
