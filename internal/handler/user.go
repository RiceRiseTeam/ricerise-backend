package handler

import (
	"ricerise/internal/apperror"
	"ricerise/internal/dto"
	"ricerise/internal/dto/request"
	"ricerise/internal/dto/response"
	"ricerise/internal/middleware"
	"ricerise/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type UserHandler struct {
	userService    *service.UserService
	authMiddleware *middleware.AuthMiddleware
}

func (h UserHandler) RegisterRouters(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	auth.POST("/register", dto.RouteWithDto(h.Register))
	auth.POST("/login", dto.RouteWithDto(h.Login))
	auth.POST("/refresh", h.Refresh)

	api := router.Group("/user")
	api.Use(h.authMiddleware.Handle)
	api.POST("/logout", dto.RouteWithDto(h.Logout))
}

func (h UserHandler) Register(context *gin.Context, req request.UserRegisterRequest) any {
	if err := h.userService.Register(context, &req); err == nil {
		return dto.Success(nil)
	}
	return nil
}

func (h UserHandler) Login(context *gin.Context, req request.UserLoginRequest) any {
	token, err := h.userService.Login(context, &req)
	if err != nil {
		return nil
	}
	return dto.Success(&response.UserLoginResponse{
		AccessToken: token,
	})
}

func (h UserHandler) Logout(context *gin.Context, _ dto.Empty) any {
	h.userService.Logout(context)
	return dto.Success(nil)
}

func (h UserHandler) Refresh(context *gin.Context) {
	oldToken, err := context.Cookie("refresh_token")
	if err != nil {
		_ = context.Error(apperror.NoRefreshTokenError)
		return
	}
	err = h.userService.Refresh(context, oldToken)
	if err != nil {
		_ = context.Error(err)
	}
}

func NewUserHandler(injector do.Injector) (*UserHandler, error) {
	authMiddleware := do.MustInvoke[*middleware.AuthMiddleware](injector)
	userService := do.MustInvoke[*service.UserService](injector)
	return &UserHandler{
		userService:    userService,
		authMiddleware: authMiddleware,
	}, nil
}
