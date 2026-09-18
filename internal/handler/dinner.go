package handler

import (
	"crypto/rand"
	"encoding/hex"
	"ricerise/internal/apperror"
	"ricerise/internal/dto"
	"ricerise/internal/dto/request"
	"ricerise/internal/dto/response"
	"ricerise/internal/middleware"
	"ricerise/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type DinnerHandler struct {
	appConfig      *config.AppConfig
}

func (h UserHandler) RegisterRouters(router *gin.RouterGroup) {
	api := router.Group("/dinner")
	api.POST("/list", dto.RouteWithDto(h.list))
}

func (h DinnerHandler) List(context *gin.Context) any{
	if err := h.repo.FindAllByA(context Request.Context(), &model.DinnerModel{
		Status: 0,
	}
	); err == nil {
		return dto.Success(nil)
	}
	return nil

}

func (h DinnerHandler) LikeFind(context *gin.context) any{
	if err := h.repo.FindAllByA(context Request.Context(), &model.DinnerModel{
		Location: request.LocationName,
	}
	); err == nil {
		return dto.Success(nil)
	}
	return nil
}