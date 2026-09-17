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

func (mw ErrorMiddleware) handle(context *gin.Context) {
	context.Next()
	if len(context.Errors) == 0 {
		return
	}

	if context.Writer.Written() {
		return
	}

	var lastErr = context.Errors.Last()
	if lastErr == nil {
		return
	}

	var appError *apperror.AppError
	if errors.As(lastErr.Err, &appError) {
		context.JSON(200, dto.Error(appError))
		return
	}

	var vError validator.ValidationErrors
	if errors.As(lastErr.Err, &vError) {
		context.JSON(200, dto.Error(apperror.ValidationError))
		return
	}

	logger.Error("Error when handling request: %w", lastErr)
	context.JSON(200, dto.Error(apperror.InternalServerError))
}

func NewErrorMiddleware(_ do.Injector) (*ErrorMiddleware, error) {
	return &ErrorMiddleware{}, nil
}
