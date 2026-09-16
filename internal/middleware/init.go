package middleware

import "github.com/gin-gonic/gin"

type Middleware interface {
	Handle(context *gin.Context)
}
