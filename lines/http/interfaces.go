package http

import "github.com/gin-gonic/gin"

type Route func(ctx *gin.Context)

type HttpEngine interface {
	Run(addr ...string) error
	GET(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	POST(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	PUT(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	DELETE(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	OPTIONS(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
}

type HttpError struct {
	Message []string `json:"message"`
}
