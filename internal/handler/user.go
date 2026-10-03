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
	"time"

	sse "github.com/dan-sherwin/go-sse"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type UserHandler struct {
	userService    *service.UserService
	authMiddleware *middleware.AuthMiddleware
}

func (h UserHandler) RegisterRouters(router *gin.RouterGroup) {
	router.GET("/sse", h.authMiddleware.CreateHandler(h.authMiddleware.UserLevel), h.SSE)

	auth := router.Group("/auth")
	auth.POST("/register", dto.RouteWithDto(h.Register))
	auth.POST("/login", dto.RouteWithDto(h.Login))
	auth.POST("/refresh", dto.RouteWithDto(h.Refresh))

	api := router.Group("/user")
	api.Use(h.authMiddleware.CreateHandler(h.authMiddleware.UserLevel))
	api.POST("/logout", dto.RouteWithDto(h.Logout))
}

// Register
// @Summary      用户注册
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request body request.UserRegisterRequest true "用户注册请求体"
// @Success      200   {object}  dto.CommonResponse  "用户注册成功"
// @Router       /auth/register [post]
func (h UserHandler) Register(ctx *gin.Context, req request.UserRegisterRequest) (any, error) {
	err := h.userService.Register(ctx, &req)
	if err != nil {
		return nil, err
	}
	return dto.Success(nil), nil
}

// Login
// @Summary      用户登录
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request body request.UserLoginRequest true "用户登录请求体"
// @Success      200   {object}  dto.CommonResponse{data=response.UserLoginResponse}  "用户登录成功"
// @Router       /auth/login [post]
func (h UserHandler) Login(ctx *gin.Context, req request.UserLoginRequest) (any, error) {
	token, err := h.userService.Login(ctx, &req)
	if err != nil {
		return nil, err
	}
	return dto.Success(&response.UserLoginResponse{
		AccessToken: token,
	}), nil
}

func (h UserHandler) Logout(ctx *gin.Context, _ dto.EmptyDto) (any, error) {
	h.userService.Logout(ctx)
	return dto.Success(nil), nil
}

// Refresh
// @Summary      刷新访问令牌
// @Tags         user
// @Router       /auth/refresh [post]
func (h UserHandler) Refresh(ctx *gin.Context, _ dto.EmptyDto) (any, error) {
	oldToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		return nil, apperror.NoRefreshTokenError
	}
	token, err := h.userService.Refresh(ctx, oldToken)
	if err != nil {
		return nil, err
	}

	return dto.Success(response.UserLoginResponse{
		AccessToken: token,
	}), nil
}

// SSE
// @Summary      服务器发送事件
// @Tags         user
// @Router       /sse [get]
func (h UserHandler) SSE(ctx *gin.Context) {
	userInfo := h.authMiddleware.GetUserInfo(ctx)
	if userInfo == nil {
		_ = ctx.Error(apperror.NoAccessTokenError)
		ctx.AbortWithStatus(401)
		return
	}

	// 生成session 随机字符串
	session := make([]byte, 16)
	if _, err := rand.Read(session); err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	sessionId := hex.EncodeToString(session)
	sse.NewSessionWithUID(ctx, sessionId, userInfo.Username)
	time.AfterFunc(5*time.Minute, func() {
		sse.ShutdownBySessionID(sessionId)
	})
}

func NewUserHandler(injector do.Injector) (*UserHandler, error) {
	authMiddleware := do.MustInvoke[*middleware.AuthMiddleware](injector)
	userService := do.MustInvoke[*service.UserService](injector)
	return &UserHandler{
		userService:    userService,
		authMiddleware: authMiddleware,
	}, nil
}
