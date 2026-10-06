package handler

import (
	"ricerise/internal/config"
	"ricerise/internal/middleware"
	"ricerise/internal/service"

	"github.com/samber/do/v2"
)

type AgentHandler struct {
	appConfig    *config.AppConfig
	auth         *middleware.AuthMiddleware
	agentService *service.AgentService
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
