package http

import "github.com/gin-gonic/gin"

type Route func(ctx *gin.Context)
