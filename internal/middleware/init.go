package middleware

import "github.com/gin-gonic/gin"

type Middleware interface {
	CreateHandler(args ...any) gin.HandlerFunc
}
