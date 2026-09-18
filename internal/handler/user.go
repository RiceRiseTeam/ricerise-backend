package handler

import (
	"crypto/rand"
	"encoding/hex"
	"ricerise/internal/apperror"
	"ricerise/internal/dto"
	"ricerise/internal/dto/request"
	"ricerise/internal/dto/response"
	"ricerise/internal/middleware"
	"ricerise/internal/service"

	sse "github.com/dan-sherwin/go-sse"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type UserHandler struct {
	userService    *service.UserService
	authMiddleware *middleware.AuthMiddleware
}

func (h UserHandler) RegisterRouters(router *gin.RouterGroup) {
	router.GET("/sse", h.SSE)

	auth := router.Group("/auth")
	auth.POST("/register", dto.RouteWithDto(h.Register))
	auth.POST("/login", dto.RouteWithDto(h.Login))
	auth.POST("/refresh", h.Refresh)

	api := router.Group("/user")
	api.Use(h.authMiddleware.CreateHandler(h.authMiddleware.UserLevel))
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

func (h UserHandler) SSE(context *gin.Context) {
	accessToken, err := context.Cookie("access_token")
	if err != nil {
		_ = context.Error(apperror.NoAccessTokenError)
		return
	}
	info := h.authMiddleware.ParseAccessToken(accessToken)
	if info == nil {
		_ = context.Error(apperror.NoAccessTokenError)
		return
	}

	// 生成session 随机字符串
	session := make([]byte, 16)
	if _, err = rand.Read(session); err != nil {
		_ = context.Error(err)
		return
	}

	sse.NewSessionWithUID(context, hex.EncodeToString(session), info.Username)
}

func NewUserHandler(injector do.Injector) (*UserHandler, error) {
	authMiddleware := do.MustInvoke[*middleware.AuthMiddleware](injector)
	userService := do.MustInvoke[*service.UserService](injector)
	return &UserHandler{
		userService:    userService,
		authMiddleware: authMiddleware,
	}, nil
}
