package middleware

import (
	"errors"
	"ricerise/internal/apperror"
	"ricerise/internal/dto"
	"ricerise/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/samber/do/v2"
)

type ErrorMiddleware struct {
}

func (mw ErrorMiddleware) CreateHandler() gin.HandlerFunc {
	return mw.handle
}

func (mw ErrorMiddleware) handle(ctx *gin.Context) {
	ctx.Next()
	if len(ctx.Errors) == 0 {
		return
	}

	if ctx.Writer.Written() {
		return
	}

	var lastErr = ctx.Errors.Last()
	if lastErr == nil {
		return
	}

	var appError *apperror.AppError
	if errors.As(lastErr.Err, &appError) {
		ctx.JSON(200, dto.Error(appError))
		return
	}

	var vError validator.ValidationErrors
	if errors.As(lastErr.Err, &vError) {
		ctx.JSON(200, dto.Error(apperror.ValidationError))
		return
	}

	logger.Error("Error when handling request: %w", lastErr)
	ctx.JSON(200, dto.Error(apperror.InternalServerError))
}

func NewErrorMiddleware(_ do.Injector) (*ErrorMiddleware, error) {
	return &ErrorMiddleware{}, nil
}
