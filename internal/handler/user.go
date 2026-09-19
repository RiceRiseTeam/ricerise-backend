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

func (h UserHandler) Register(ctx *gin.Context, req request.UserRegisterRequest) any {
	if err := h.userService.Register(ctx, &req); err == nil {
		return dto.Success(nil)
	}
	return nil
}

func (h UserHandler) Login(ctx *gin.Context, req request.UserLoginRequest) any {
	token, err := h.userService.Login(ctx, &req)
	if err != nil {
		return nil
	}
	return dto.Success(&response.UserLoginResponse{
		AccessToken: token,
	})
}

func (h UserHandler) Logout(ctx *gin.Context, _ dto.Empty) any {
	h.userService.Logout(ctx)
	return dto.Success(nil)
}

func (h UserHandler) Refresh(ctx *gin.Context) {
	oldToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		_ = ctx.Error(apperror.NoRefreshTokenError)
		return
	}
	err = h.userService.Refresh(ctx, oldToken)
	if err != nil {
		_ = ctx.Error(err)
	}
}

func (h UserHandler) SSE(ctx *gin.Context) {
	accessToken, err := ctx.Cookie("access_token")
	if err != nil {
		_ = ctx.Error(apperror.NoAccessTokenError)
		return
	}
	info := h.authMiddleware.ParseAccessToken(accessToken)
	if info == nil {
		_ = ctx.Error(apperror.NoAccessTokenError)
		return
	}

	// 生成session 随机字符串
	session := make([]byte, 16)
	if _, err = rand.Read(session); err != nil {
		_ = ctx.Error(err)
		return
	}

	sse.NewSessionWithUID(ctx, hex.EncodeToString(session), info.Username)
}

func NewUserHandler(injector do.Injector) (*UserHandler, error) {
	authMiddleware := do.MustInvoke[*middleware.AuthMiddleware](injector)
	userService := do.MustInvoke[*service.UserService](injector)
	return &UserHandler{
		userService:    userService,
		authMiddleware: authMiddleware,
	}, nil
}
