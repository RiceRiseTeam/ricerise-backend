package handler

import (
	"ricerise/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type UserHandler struct {
	serv *service.UserService
}

func (h UserHandler) RegisterRouters(router *gin.RouterGroup) {
	router.Group("/user")
}

func NewUserHandler(injector do.Injector) (*UserHandler, error) {
	serv := do.MustInvoke[*service.UserService](injector)
	return &UserHandler{serv: serv}, nil
}
