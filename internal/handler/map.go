package handler

import (
	"ricerise/internal/config"
	"ricerise/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type MapHandler struct {
	appConfig  *config.AppConfig
	mapService *service.MapService
}

func (m MapHandler) RegisterRouters(router *gin.RouterGroup) {
}

func NewMapHandler(injector do.Injector) (*MapHandler, error) {
	mapService := do.MustInvoke[*service.MapService](injector)
	appConfig := do.MustInvoke[*config.AppConfig](injector)
	return &MapHandler{
		appConfig:  appConfig,
		mapService: mapService,
	}, nil
}
