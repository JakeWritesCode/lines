package http

import "github.com/gin-gonic/gin"

type Route func(ctx *gin.Context)

type HttpEngine interface {
	Run(addr ...string) error
}

type HttpError struct {
	Message []string `json:"message"`
}
