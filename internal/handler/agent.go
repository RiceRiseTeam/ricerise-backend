package handler

import (
	"ricerise/internal/apperror"
	"ricerise/internal/config"
	"ricerise/internal/dto"
	"ricerise/internal/dto/request"
	"ricerise/internal/dto/response"
	"ricerise/internal/middleware"
	"ricerise/internal/service"
	"strconv"

	"github.com/dan-sherwin/go-sse"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
)

type AgentHandler struct {
	appConfig    *config.AppConfig
	auth         *middleware.AuthMiddleware
	agentService *service.AgentService
}

func (a AgentHandler) RegisterRouters(router *gin.RouterGroup) {
	router.GET("/chat/:id/stream", a.ChatSSE)

	api := router.Group("chat")
	api.Use(a.auth.CreateHandler(a.auth.UserLevel))
	api.GET("/:id/messages", dto.RouteWithDto(a.GetMessagesHistory))
	api.POST("/:id/messages", dto.RouteWithDto(a.Chat))
	api.GET("/session", dto.RouteWithDto(a.GetSession))
}

// Chat
// @Summary      向agent发送消息
// @Tags         agent
// @Param        id path int true "sessionId"
// @Param        request body request.AgentChatRequest true "Agent消息请求体"
// @Success      200   {object}  dto.CommonResponse  "发送消息成功"
// @Router       /chat/:id/messages [post]
func (a AgentHandler) Chat(ctx *gin.Context, request request.AgentChatRequest) (any, error) {
	userInfo := a.auth.GetUserInfo(ctx)
	if userInfo == nil {
		return nil, apperror.NoAccessTokenError
	}
	sessionId, err := dto.GetUrlID(ctx)
	if err != nil {
		return nil, apperror.AccessNoFoundError
	}

	err = a.agentService.Chat(ctx, userInfo.UserId, userInfo.Username, sessionId, request.Input)
	if err != nil {
		return nil, err
	}

	return dto.Success(nil), nil
}

// GetMessagesHistory
// @Summary      获取对话历史
// @Tags         agent
// @Param        id path int true "sessionId"
// @Success      200   {object}  dto.CommonResponse{data=[]dto.ChatMessageDto}  "获取对话历史成功"
// @Router       /chat/:id/messages [get]
func (a AgentHandler) GetMessagesHistory(ctx *gin.Context, _ dto.EmptyDto) (any, error) {
	userInfo := a.auth.GetUserInfo(ctx)
	if userInfo == nil {
		return nil, apperror.NoAccessTokenError
	}
	sessionId, err := dto.GetUrlID(ctx)
	if err != nil {
		return nil, apperror.AccessNoFoundError
	}

	result, err := a.agentService.GetMessageHistory(ctx, userInfo.UserId, sessionId)
	if err != nil {
		return nil, err
	}

	return dto.Success(result), nil
}

// GetSession
// @Summary      新对话获取新的SessionId
// @Tags         agent
// @Success      200   {object}  dto.CommonResponse{data=response.AgentSessionResponse}  "获取SessionId成功"
// @Router       /chat/session [get]
func (a AgentHandler) GetSession(ctx *gin.Context, _ dto.EmptyDto) (any, error) {
	userInfo := a.auth.GetUserInfo(ctx)
	if userInfo == nil {
		return nil, apperror.NoAccessTokenError
	}
	sessionId, err := a.agentService.GetSessionId(ctx, userInfo.UserId)
	if err != nil {
		return nil, err
	}

	return dto.Success(response.AgentSessionResponse{SessionId: sessionId}), nil
}

// ChatSSE
// @Summary      服务器发送Agent 返回的事件 每次agent发送完消息后都会关闭
// @Tags         agent
// @Param        id path int true "sessionId"
// @Router       /chat/:id/stream [get]
func (a AgentHandler) ChatSSE(ctx *gin.Context) {
	userInfo := a.auth.GetUserInfo(ctx)
	if userInfo == nil {
		_ = ctx.Error(apperror.NoAccessTokenError)
		ctx.AbortWithStatus(401)
		return
	}

	sessionId, err := dto.GetUrlID(ctx)
	if err != nil {
		_ = ctx.Error(apperror.AccessNoFoundError)
		ctx.AbortWithStatus(404)
		return
	}

	sse.NewSessionWithUID(ctx, uuid.New(), "agent-"+strconv.FormatUint(sessionId, 10)+userInfo.Username)
}

func NewAgentHandler(injector do.Injector) (*AgentHandler, error) {
	authMiddleware := do.MustInvoke[*middleware.AuthMiddleware](injector)
	agentService := do.MustInvoke[*service.AgentService](injector)
	appConfig := do.MustInvoke[*config.AppConfig](injector)
	return &AgentHandler{
		appConfig:    appConfig,
		auth:         authMiddleware,
		agentService: agentService,
	}, nil
}
