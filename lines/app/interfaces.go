package app

import (
	"github.com/gin-gonic/gin"
	"lines/internal"
)

// App is the interface that all apps must implement.
type App interface {
	// Initialise is called to initialise the app.
	Initialise(mainConfig *internal.MainConfig) error
	// RegisterHTTPRoutes is called to register the app's HTTP routes.
	RegisterHTTPRoutes(engine *gin.Engine)
	// RegisterGRPCServices is called to register the app's gRPC services.
	RegisterGRPCServices() error
}
