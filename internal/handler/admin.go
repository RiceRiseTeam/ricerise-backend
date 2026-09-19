package handler

import (
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

	api.GET("/comments", dto.RouteQueryWithDto(a.GetComments)) // 获取需要审核的Comment和 Location
	api.GET("/locations", dto.RouteQueryWithDto(a.GetLocations))
	api.GET("/status", dto.RouteQueryWithDto(a.GetStatus))
}

func (a AdminHandler) GetStatus(ctx *gin.Context, _ dto.Empty) any {
	return a.adminService.GetAppStatus(ctx)
}

func (a AdminHandler) GetComments(ctx *gin.Context, query query.AdminPageQuery) any {
	return nil
}

func (a AdminHandler) GetLocations(ctx *gin.Context, query query.AdminPageQuery) any {
	data, hasNext, err := a.adminService.GetLocationsReviewList(ctx, query)
	if err != nil {
		return nil
	}
	return dto.Success(&response.AdminLocationReviewList{
		PageSize:  query.PageSize,
		HasNext:   hasNext,
		Locations: dto.Map(data, dto.NewLocationDto),
	})
}

func (a AdminHandler) ReviewComment(ctx *gin.Context) any {
	return nil
}

func (a AdminHandler) ReviewLocation(ctx *gin.Context) any {
	return nil
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
