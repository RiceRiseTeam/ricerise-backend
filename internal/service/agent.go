package service

import "github.com/samber/do/v2"

type AgentService struct {
}

func NewAgentService(injector do.Injector) (*AgentService, error) {
	return &AgentService{}, nil
}
