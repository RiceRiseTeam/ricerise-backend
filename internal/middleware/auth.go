package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type AuthMiddleware struct {
}

func (a AuthMiddleware) Handle(context *gin.Context) {

}

func NewAuthMiddleware(injector do.Injector) (*AuthMiddleware, error) {
	return &AuthMiddleware{}, nil
}
