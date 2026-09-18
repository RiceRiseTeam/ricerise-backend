package service

import (
	"ricerise/internal/config"
	"ricerise/internal/repository"

	"github.com/samber/do/v2"
)

type AdminService struct {
	userRepository *repository.UserRepository
	appConfig      *config.AppConfig
}

func NewAdminService(injector do.Injector) (*AdminService, error) {
	userRepository := do.MustInvoke[*repository.UserRepository](injector)
	appConfig := do.MustInvoke[*config.AppConfig](injector)
	return &AdminService{
		userRepository: userRepository,
		appConfig:      appConfig,
	}, nil
}
