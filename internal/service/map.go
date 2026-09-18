package service

import (
	"ricerise/internal/config"

	"github.com/samber/do/v2"
)

type MapService struct {
	appConfig *config.AppConfig
}

func NewMapService(injector do.Injector) (*MapService, error) {
	appConfig := do.MustInvoke[*config.AppConfig](injector)
	return &MapService{
		appConfig: appConfig,
	}, nil
}
