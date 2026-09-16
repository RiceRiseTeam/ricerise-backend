package handler

import (
	"ricerise/internal/dto"
	"ricerise/internal/dto/request"
	"ricerise/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type UserHandler struct {
	serv *service.UserService
}

func (h UserHandler) RegisterRouters(router *gin.RouterGroup) {
	api := router.Group("/user")
	api.POST("/register", dto.RouteWithDto(h.Register))
}

func (h UserHandler) Register(context *gin.Context, req request.UserRegisterRequest) any {
	h.serv.RegisterNew(&req)
	return dto.Success(nil)
}

func NewUserHandler(injector do.Injector) (*UserHandler, error) {
	serv := do.MustInvoke[*service.UserService](injector)
	return &UserHandler{serv: serv}, nil
}
