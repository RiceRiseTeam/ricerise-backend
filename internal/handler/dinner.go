package handler

import (
	"ricerise/internal/dto"
	"ricerise/internal/dto/request"
	"ricerise/internal/dto/response"
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
	api := router.Group("/dinner")
	api.GET("/list", dto.RouteJsonWithDto(h.ListDinner))
	api.POST("/likefind", dto.RouteJsonWithDto(h.LikeFind))
	api.GET("/:id/participants", h.ListParticipant)
	api.POST("/newparticipate", dto.RouteJsonWithDto(h.NewParticipate))
}

func (h DinnerHandler) ListDinner(ctx *gin.Context, _ dto.Empty) any {
	result, err := h.dinnerService.List(ctx.Request.Context())
	if err != nil {
		_ = ctx.Error(err)
		return nil
	}
	return dto.Success(&response.DinnerFindResponse{Result: *result})
}

func (h DinnerHandler) LikeFind(ctx *gin.Context, req request.DinnerLikeFindRequest) any {
	result, err := h.dinnerService.LikeFind(ctx.Request.Context(), req.LocationName)
	if err != nil {
		_ = ctx.Error(err)
		return nil
	}
	return dto.Success(&response.DinnerFindResponse{Result: *result})
}

func (h DinnerHandler) ListParticipant(ctx *gin.Context) {
	var uri struct {
		ID uint64 `uri:"id" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&uri); err != nil {
		_ = ctx.Error(err)
		return
	}
	result, err := h.dinnerService.ListParticipant(ctx.Request.Context(), uri.ID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(200, dto.Success(&response.ParticipantFindResponse{Result: *result}))
}

func (h DinnerHandler) NewParticipate(ctx *gin.Context, req request.NewParticipateRequest) any {
	info := h.auth.GetUserInfo(ctx)
	if info == nil {
		return nil
	}
	if err := h.dinnerService.NewParticipate(ctx.Request.Context(), info.UserId, req.DinnerID); err != nil {
		_ = ctx.Error(err)
		return nil
	}
	return dto.Success(nil)
}

func NewDinnerHandler(injector do.Injector) (*DinnerHandler, error) {
	return &DinnerHandler{
		dinnerService: do.MustInvoke[*service.DinnerService](injector),
		auth:          do.MustInvoke[*middleware.AuthMiddleware](injector),
	}, nil
}
