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

func RouteJsonWithDto[T any](input func(ctx *gin.Context, dto T) any) func(ctx *gin.Context) {
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

func RouteQueryWithDto[T any](input func(ctx *gin.Context, dto T) any) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req T
		var result any
		switch any(req).(type) {
		case Empty:
			result = input(ctx, req)
		default:
			err := ctx.ShouldBindQuery(&req)
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

func GetUrlID(ctx *gin.Context) (uint64, error) {
	var uri struct {
		ID uint64 `uri:"id" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&uri); err != nil {
		return 0, err
	}

	return uri.ID, nil
}

// Map Java Collection.stream().map() 类似物 用于 model -> dto 之间的装包
func Map[T, U any](array []T, mapper func(T) U) []U {
	result := make([]U, 0, len(array))
	for _, v := range array {
		result = append(result, mapper(v))
	}
	return result
}
