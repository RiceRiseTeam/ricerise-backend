package handler

import (
	"ricerise/internal/config"
	"ricerise/internal/dto"
	"ricerise/internal/dto/request"
	"ricerise/internal/middleware"
	"ricerise/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type MapHandler struct {
	appConfig      *config.AppConfig
	mapService     *service.MapService
	authMiddleware *middleware.AuthMiddleware
}

func (m MapHandler) RegisterRouters(router *gin.RouterGroup) {
	api := router.Group("/map")
	api.Use(m.authMiddleware.CreateHandler(m.authMiddleware.UserLevel))

	api.POST("/comments", dto.RouteWithDto(m.UploadComment))
	api.POST("/locations", dto.RouteWithDto(m.UploadLocation))
}

func (m MapHandler) UploadComment(ctx *gin.Context, request request.UploadCommentRequest) (any, error) {
	result, err := m.mapService.UploadComment(ctx, request)
	if err != nil {
		return nil, err
	}
	return dto.Created(result), nil
}

func (m MapHandler) UploadLocation(ctx *gin.Context, request request.UploadLocationRequest) (any, error) {
	result, err := m.mapService.UploadLocation(ctx, request)
	if err != nil {
		return nil, nil
	}
	return dto.Created(result), nil
}

func NewMapHandler(injector do.Injector) (*MapHandler, error) {
	authMiddleware := do.MustInvoke[*middleware.AuthMiddleware](injector)
	mapService := do.MustInvoke[*service.MapService](injector)
	appConfig := do.MustInvoke[*config.AppConfig](injector)
	return &MapHandler{
		appConfig:      appConfig,
		mapService:     mapService,
		authMiddleware: authMiddleware,
	}, nil
}
