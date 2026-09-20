package service

import (
	"errors"
	"ricerise/internal/apperror"
	"ricerise/internal/config"
	"ricerise/internal/dto"
	"ricerise/internal/dto/request"
	"ricerise/internal/middleware"
	"ricerise/internal/model"
	"ricerise/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type MapService struct {
	appConfig          *config.AppConfig
	locationRepository *repository.LocationRepository
	commentRepository  *repository.CommentRepository
	authMiddleware     *middleware.AuthMiddleware
}

func (m MapService) UploadComment(ctx *gin.Context, request request.UploadCommentRequest) *dto.CommentDto {
	userInfo := m.authMiddleware.GetUserInfo(ctx)
	if userInfo == nil {
		_ = ctx.Error(apperror.NoAccessTokenError)
		return nil
	}
	goContext := ctx.Request.Context()
	location, err := m.locationRepository.FindById(goContext, request.LocationId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = ctx.Error(apperror.AccessNoFoundError)
			return nil
		}
		_ = ctx.Error(err)
		return nil
	}

	newComment := &model.CommentModel{
		Rating:     request.Rating,
		Content:    request.Content,
		UserId:     userInfo.UserId,
		LocationId: location.ID,
	}
	err = m.commentRepository.Create(goContext, newComment)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			_ = ctx.Error(apperror.NameConflictError)
			return nil
		}
		_ = ctx.Error(err)
		return nil
	}

	return dto.NewCommentDto(newComment)
}

func (m MapService) UploadLocation(ctx *gin.Context, request request.UploadLocationRequest) *dto.LocationDto {
	userInfo := m.authMiddleware.GetUserInfo(ctx)
	if userInfo == nil {
		_ = ctx.Error(apperror.NoAccessTokenError)
		return nil
	}

	newLocation := &model.LocationModel{
		Longitude:   request.Longitude,
		Latitude:    request.Latitude,
		Name:        request.Name,
		Description: request.Description,
	}

	err := m.locationRepository.Create(ctx.Request.Context(), newLocation)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			_ = ctx.Error(apperror.NameConflictError)
			return nil
		}
		_ = ctx.Error(err)
		return nil
	}

	return dto.NewLocationDto(newLocation)
}

func NewMapService(injector do.Injector) (*MapService, error) {
	commentRepository := do.MustInvoke[*repository.CommentRepository](injector)
	locationRepository := do.MustInvoke[*repository.LocationRepository](injector)
	auth := do.MustInvoke[*middleware.AuthMiddleware](injector)
	appConfig := do.MustInvoke[*config.AppConfig](injector)
	return &MapService{
		appConfig:          appConfig,
		authMiddleware:     auth,
		locationRepository: locationRepository,
		commentRepository:  commentRepository,
	}, nil
}
