package handler

import (
	"ricerise/internal/apperror"
	"ricerise/internal/config"
	"ricerise/internal/dto"
	"ricerise/internal/dto/query"
	"ricerise/internal/dto/response"
	"ricerise/internal/middleware"
	"ricerise/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type AdminHandler struct {
	appConfig      *config.AppConfig
	authMiddleware *middleware.AuthMiddleware
	adminService   *service.AdminService
}

func (a AdminHandler) RegisterRouters(router *gin.RouterGroup) {
	api := router.Group("/admin")
	api.Use(a.authMiddleware.CreateHandler(a.authMiddleware.AdminLevel))

	api.GET("/comments", dto.RouteWithDto(a.GetComments)) // 获取需要审核的Comment和 Location
	api.GET("/locations", dto.RouteWithDto(a.GetLocations))
	api.GET("/status", dto.RouteWithDto(a.GetStatus))
	api.PATCH("/comments/:id", dto.RouteWithDto(a.ReviewComment))
	api.PATCH("/locations/:id", dto.RouteWithDto(a.ReviewLocation))
}

func (a AdminHandler) GetStatus(ctx *gin.Context, _ dto.EmptyDto) (any, error) {
	return dto.Success(a.adminService.GetAppStatus(ctx)), nil
}

func (a AdminHandler) GetComments(ctx *gin.Context, query query.AdminPageQuery) (any, error) {
	data, hasNext, err := a.adminService.GetCommentReviewList(ctx, query)
	if err != nil {
		return nil, err
	}
	return dto.Success(&response.AdminCommentReviewList{
		PageSize: query.PageSize,
		HasNext:  hasNext,
		Comments: dto.Map(data, dto.NewCommentDto),
	}), nil
}

func (a AdminHandler) GetLocations(ctx *gin.Context, query query.AdminPageQuery) (any, error) {
	data, hasNext, err := a.adminService.GetLocationReviewList(ctx, query)
	if err != nil {
		return nil, err
	}
	return dto.Success(&response.AdminLocationReviewList{
		PageSize:  query.PageSize,
		HasNext:   hasNext,
		Locations: dto.Map(data, dto.NewLocationDto),
	}), nil
}

func (a AdminHandler) ReviewComment(ctx *gin.Context, query query.AdminReviewQuery) (any, error) {
	commentId, err := dto.GetUrlID(ctx)
	if err != nil {
		return dto.Error(apperror.ValidationError), nil
	}
	err = a.adminService.ReviewLocation(ctx, commentId, query.Pass)
	if err != nil {
		return nil, err
	}
	return dto.Success(nil), nil
}

func (a AdminHandler) ReviewLocation(ctx *gin.Context, query query.AdminReviewQuery) (any, error) {
	locationId, err := dto.GetUrlID(ctx)
	if err != nil {
		return dto.Error(apperror.ValidationError), nil
	}
	err = a.adminService.ReviewLocation(ctx, locationId, query.Pass)
	if err != nil {
		return nil, err
	}
	return dto.Success(nil), nil
}

func NewAdminHandler(injector do.Injector) (*AdminHandler, error) {
	authMiddleware := do.MustInvoke[*middleware.AuthMiddleware](injector)
	adminService := do.MustInvoke[*service.AdminService](injector)
	appConfig := do.MustInvoke[*config.AppConfig](injector)
	return &AdminHandler{
		appConfig:      appConfig,
		authMiddleware: authMiddleware,
		adminService:   adminService,
	}, nil
}
