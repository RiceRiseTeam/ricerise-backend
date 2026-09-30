package handler

import (
	"ricerise/internal/apperror"
	"ricerise/internal/config"
	"ricerise/internal/dto"
	"ricerise/internal/dto/querydto"
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

// GetStatus
// @Summary      获取应用状态
// @Tags         admin
// @Accept       json
// @Produce      json
// @Success      200   {object}  dto.CommonResponse{data=response.AdminStatusResponse}  "获取应用状态成功"
// @Router       /admin/status [get]
func (a AdminHandler) GetStatus(ctx *gin.Context, _ dto.EmptyDto) (any, error) {
	return dto.Success(a.adminService.GetAppStatus(ctx)), nil
}

// GetComments
// @Summary      获取需要审核的评论列表
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        query query querydto.AdminPageQuery true "分页查询参数"
// @Success      200   {object}  dto.CommonResponse{data=response.AdminCommentReviewList}  "获取需要审核的评论列表成功"
// @Router       /admin/comments [get]
func (a AdminHandler) GetComments(ctx *gin.Context, query querydto.AdminPageQuery) (any, error) {
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

// GetLocations
// @Summary      获取需要审核的地点列表
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        query query querydto.AdminPageQuery true "分页查询参数"
// @Success      200   {object}  dto.CommonResponse{data=response.AdminLocationReviewList}  "获取需要审核的地点列表成功"
// @Router       /admin/locations [get]
func (a AdminHandler) GetLocations(ctx *gin.Context, query querydto.AdminPageQuery) (any, error) {
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

// ReviewComment
// @Summary      审核评论
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id path int true "评论ID"
// @Param        query query querydto.AdminReviewQuery true "审核参数"
// @Success      200   {object}  dto.CommonResponse  "审核评论成功"
// @Router       /admin/comments/:id [post]
func (a AdminHandler) ReviewComment(ctx *gin.Context, query querydto.AdminReviewQuery) (any, error) {
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

// ReviewLocation
// @Summary      审核地点
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id path int true "地点ID"
// @Param        query query querydto.AdminReviewQuery true "审核参数"
// @Success      200   {object}  dto.CommonResponse  "审核地点成功"
// @Router       /admin/locations/:id [post]
func (a AdminHandler) ReviewLocation(ctx *gin.Context, query querydto.AdminReviewQuery) (any, error) {
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
