package handler

import (
	"ricerise/internal/apperror"
	"ricerise/internal/dto"
	"ricerise/internal/dto/request"
	"ricerise/internal/middleware"
	"ricerise/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type DinnerHandler struct {
	dinnerService *service.DinnerService
	auth          *middleware.AuthMiddleware
}

func (h DinnerHandler) RegisterRouters(router *gin.RouterGroup) {
	api := router.Group("/dinners")
	api.Use(h.auth.CreateHandler(h.auth.UserLevel))
	router.POST("/", dto.RouteWithDto(h.CreateDinner))

	router.POST("/:id/participants/", dto.RouteWithDto(h.JointDinner))
	router.DELETE("/:id/participants/me", dto.RouteWithDto(h.LeftDinner))
	router.PATCH("/:id/", dto.RouteWithDto(h.UpdateDinnerStatus))
	router.POST("/:id/messages", dto.RouteWithDto(h.DinnerRoomChat))
}

// DinnerRoomChat
// @Summary      饭局内聊天
// @Tags         dinner
// @Accept       json
// @Produce      json
// @Param        id path int true "饭局ID"
// @Param        request body request.DinnerRoomChatRequest true "饭局聊天请求体"
// @Success      200   {object}  dto.CommonResponse  "发送消息成功"
// @Router       /dinners/:id [patch]
func (h DinnerHandler) DinnerRoomChat(ctx *gin.Context, request request.DinnerRoomChatRequest) (any, error) {
	userInfo := h.auth.GetUserInfo(ctx)
	if userInfo == nil {
		return nil, apperror.NoAccessTokenError
	}
	dinnerId, err := dto.GetUrlID(ctx)
	if err != nil {
		return nil, err
	}
	err = h.dinnerService.SendDinnerMessage(ctx, userInfo.UserId, dinnerId, request.Message)
	if err != nil {
		return nil, err
	}

	return dto.Success(nil), nil
}

// CreateDinner
// @Summary      创建饭局
// @Tags         dinner
// @Accept       json
// @Produce      json
// @Param        request body request.CreateDinnerRequest true "创建饭局请求体"
// @Success      200   {object}  dto.CommonResponse{data=dto.DinnerDto}  "修改状态成功"
// @Router       /dinners/:id [patch]
func (h DinnerHandler) CreateDinner(ctx *gin.Context, request request.CreateDinnerRequest) (any, error) {
	userInfo := h.auth.GetUserInfo(ctx)
	if userInfo == nil {
		return nil, apperror.NoAccessTokenError
	}

	dinner, err := h.dinnerService.CreateDinner(ctx, userInfo.UserId, request)
	if err != nil {
		return nil, err
	}
	return dto.Created(dinner), nil
}

// JointDinner
// @Summary      加入约饭
// @Tags         dinner
// @Accept       json
// @Produce      json
// @Param        id path int true "饭局ID"
// @Success      200   {object}  dto.CommonResponse{data=dto.DinnerDto}  "加入饭局成功"
// @Router       /dinners/:id/participants/ [post]
func (h DinnerHandler) JointDinner(ctx *gin.Context, _ dto.EmptyDto) (any, error) {
	dinnerId, err := dto.GetUrlID(ctx)
	if err != nil {
		return nil, err
	}
	userInfo := h.auth.GetUserInfo(ctx)
	if userInfo == nil {
		return nil, apperror.NoAccessTokenError
	}
	result, err := h.dinnerService.JoinDinner(ctx, userInfo.UserId, dinnerId)
	if err != nil {
		return nil, err
	}

	return dto.Success(result), nil
}

// LeftDinner
// @Summary      退出约饭
// @Tags         dinner
// @Accept       json
// @Produce      json
// @Param        id path int true "饭局ID"
// @Success      200   {object}  dto.CommonResponse  "退出约饭成功"
// @Router       /dinners/:id/participants/me [delete]
func (h DinnerHandler) LeftDinner(ctx *gin.Context, _ dto.EmptyDto) (any, error) {
	dinnerId, err := dto.GetUrlID(ctx)
	if err != nil {
		return nil, err
	}
	userInfo := h.auth.GetUserInfo(ctx)
	if userInfo == nil {
		return nil, apperror.NoAccessTokenError
	}
	err = h.dinnerService.ExistDinner(ctx, userInfo.UserId, dinnerId)
	if err != nil {
		return nil, err
	}

	return dto.Success(nil), nil
}

// UpdateDinnerStatus
// @Summary      修改饭局状态
// @Tags         dinner
// @Accept       json
// @Produce      json
// @Param        id path int true "饭局ID"
// @Param        request body request.UpdateDinnerStatusRequest true "更新状态请求体"
// @Success      200   {object}  dto.CommonResponse  "修改状态成功"
// @Router       /dinners/:id [patch]
func (h DinnerHandler) UpdateDinnerStatus(ctx *gin.Context, request request.UpdateDinnerStatusRequest) (any, error) {
	dinnerId, err := dto.GetUrlID(ctx)
	if err != nil {
		return nil, err
	}
	userInfo := h.auth.GetUserInfo(ctx)
	if userInfo == nil {
		return nil, apperror.NoAccessTokenError
	}
	err = h.dinnerService.UpdateDinnerStatus(ctx, userInfo.UserId, dinnerId, request.Status)
	if err != nil {
		return nil, err
	}
	return dto.Success(nil), err
}

func NewDinnerHandler(injector do.Injector) (*DinnerHandler, error) {
	return &DinnerHandler{
		dinnerService: do.MustInvoke[*service.DinnerService](injector),
		auth:          do.MustInvoke[*middleware.AuthMiddleware](injector),
	}, nil
}
