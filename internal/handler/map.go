package handler

import (
	"ricerise/internal/apperror"
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

	api.POST("/locations", dto.RouteWithDto(m.UploadLocation))
	api.POST("/location/view", dto.RouteWithDto(m.GetLocationsInRange))
	api.POST("/location/:id/comments", dto.RouteWithDto(m.UploadComment))

	api.GET("/location/:id", dto.RouteWithDto(m.GetLocationDetail))
	api.GET("/location/:id/comments")

}

// UploadComment
// @Summary      上传评论
// @Tags         map
// @Accept       json
// @Produce      json
// @Param        id path int true "地点ID"
// @Param        request body request.UploadCommentRequest true "上传评论请求体"
// @Success      200   {object}  dto.CommonResponse  "上传评论成功"
// @Router       /map/location/:id/comments [post]
func (m MapHandler) UploadComment(ctx *gin.Context, request request.UploadCommentRequest) (any, error) {
	result, err := m.mapService.UploadComment(ctx, request)
	if err != nil {
		return nil, err
	}
	return dto.Created(result), nil
}

// UploadLocation
// @Summary      上传地点
// @Tags         map
// @Accept       json
// @Produce      json
// @Param        request body request.UploadLocationRequest true "上传地点请求体"
// @Success      200   {object}  dto.CommonResponse  "上传地点成功"
// @Router       /map/locations [post]
func (m MapHandler) UploadLocation(ctx *gin.Context, request request.UploadLocationRequest) (any, error) {
	result, err := m.mapService.UploadLocation(ctx, request)
	if err != nil {
		return nil, nil
	}
	return dto.Created(result), nil
}

// GetLocationDetail
// @Summary      获取地点详情
// @Tags         map
// @Accept       json
// @Produce      json
// @Param        id path int true "地点ID"
// @Success      200   {object}  dto.CommonResponse  "获取地点详情成功"
// @Router       /map/location/:id [get]
func (m MapHandler) GetLocationDetail(ctx *gin.Context, _ dto.EmptyDto) (any, error) {
	id, err := dto.GetUrlID(ctx)
	if err != nil {
		return nil, apperror.AccessNoFoundError
	}
	result, err := m.mapService.GetLocationDetail(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.Success(result), nil
}

// GetLocationsInRange
// @Summary      获取指定范围内的地点列表
// @Tags         map
// @Accept       json
// @Produce      json
// @Param        request body request.GetLocationsRequest true "查询参数"
// @Success      200   {object}  dto.CommonResponse  "获取指定地点列表成功"
// @Router       /map/location/view [post]
func (m MapHandler) GetLocationsInRange(ctx *gin.Context, request request.GetLocationsRequest) (any, error) {
	result, err := m.mapService.GetLocationsInRange(ctx, request)
	if err != nil {
		return nil, err
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
