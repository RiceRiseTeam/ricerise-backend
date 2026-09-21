package dto

import (
	"ricerise/internal/apperror"

	"github.com/gin-gonic/gin"
)

type EmptyDto struct{}

type CommonResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Success(data any) *CommonResponse {
	return &CommonResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

func Created(data any) *CommonResponse {
	return &CommonResponse{
		Code:    20000,
		Message: "created",
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

func RouteWithDto[T any](input func(ctx *gin.Context, dto T) any) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req T
		var result any

		if _, ok := any(req).(EmptyDto); !ok {
			if err := ctx.ShouldBind(&req); err != nil {
				_ = ctx.Error(err)
				return
			}
		}
		result = input(ctx, req)
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
