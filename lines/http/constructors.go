package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"lines/internal"
)

// CreateEngine creates a new gin engine and sorts CORS out.
func CreateEngine(config internal.MainConfig) *gin.Engine {
	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = config.CORSOrigins
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))
	return r
}
