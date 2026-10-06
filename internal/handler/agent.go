package handler

import (
	"ricerise/internal/config"
	"ricerise/internal/dto"
	"ricerise/internal/middleware"
	"ricerise/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type AgentHandler struct {
	appConfig    *config.AppConfig
	auth         *middleware.AuthMiddleware
	agentService *service.AgentService
}

func (a AgentHandler) RegisterRouters(router *gin.RouterGroup) {
	api := router.Group("chat")
	api.Use(a.auth.CreateHandler(a.auth.UserLevel))
	api.GET("/:id/message", dto.RouteWithDto(a.GetMessagesHistory))
	api.GET("/session", dto.RouteWithDto(a.GetSession))
}

func (a AgentHandler) GetMessagesHistory(ctx *gin.Context, _ dto.EmptyDto) (any, error) {
	return nil, nil
}

func (a AgentHandler) GetSession(ctx *gin.Context, _ dto.EmptyDto) (any, error) {
	return nil, nil
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
