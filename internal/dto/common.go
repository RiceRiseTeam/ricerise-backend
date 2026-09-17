package dto

import (
	"ricerise/internal/apperror"

	"github.com/gin-gonic/gin"
)

type Empty struct{}

type CommonResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(data any) *CommonResponse {
	return &CommonResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

func Error(err *apperror.AppError) *CommonResponse {
	return &CommonResponse{
		Code:    err.Code,
		Message: err.Message,
		Data:    nil,
	}
}

var (
	FatalError = &CommonResponse{
		Code:    50000,
		Message: "internal server error",
		Data:    nil,
	}

	ValidationError = &CommonResponse{
		Code:    40000,
		Message: "参数校验失败",
		Data:    nil,
	}

	AuthError = &CommonResponse{
		Code:    40100,
		Message: "未登录或令牌失效",
		Data:    nil,
	}
)

func RouteWithDto[T any](input func(ctx *gin.Context, dto T) any) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req T
		var result any
		switch any(req).(type) {
		case Empty:
			result = input(ctx, req)
		default:
			err := ctx.ShouldBindJSON(&req)
			if err != nil {
				_ = ctx.Error(err)
				return
			}
			result = input(ctx, req)
		}

		if result != nil {
			ctx.JSON(200, result)
		}
	}
}
