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

func (m MapService) DeleteComment(ctx *gin.Context, id uint64) error {
	comment, err := m.fetchLocation(ctx, id)
	if err != nil {
		return err
	}

	//TODO(linstarowo): 鉴权
	err = m.commentRepository.Delete(ctx, comment.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m MapService) DeleteLocation(ctx *gin.Context, id uint64) error {
	location, err := m.fetchLocation(ctx, id)
	if err != nil {
		return err
	}
	//TODO(linstarowo): 鉴权
	err = m.locationRepository.Delete(ctx, location.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m MapService) fetchLocation(ctx *gin.Context, id uint64) (*model.LocationModel, error) {
	location, err := m.locationRepository.FindById(ctx.Request.Context(), id)
	if err != nil {
		return nil, err
	}
	if location == nil {
		return nil, apperror.AccessNoFoundError
	}
	return location, nil
}

func (m MapService) fetchComment(ctx *gin.Context, id uint64) (*model.CommentModel, error) {
	comment, err := m.commentRepository.FindById(ctx.Request.Context(), id)
	if err != nil {
		return nil, err
	}
	if comment == nil {
		return nil, apperror.AccessNoFoundError
	}
	return comment, nil
}

func (m MapService) UploadComment(ctx *gin.Context, request request.UploadCommentRequest) (*dto.CommentDto, error) {
	userInfo := m.authMiddleware.GetUserInfo(ctx)
	if userInfo == nil {
		return nil, apperror.NoAccessTokenError
	}
	goContext := ctx.Request.Context()
	location, err := m.fetchLocation(ctx, request.LocationId)
	if location == nil {
		return nil, err
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
			return nil, apperror.NameConflictError
		}
		return nil, err
	}

	return dto.NewCommentDto(newComment), nil
}

func (m MapService) UploadLocation(ctx *gin.Context, request request.UploadLocationRequest) (*dto.LocationDto, error) {
	userInfo := m.authMiddleware.GetUserInfo(ctx)
	if userInfo == nil {
		return nil, apperror.NoAccessTokenError
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
			return nil, apperror.NameConflictError
		}
		return nil, err
	}

	return dto.NewLocationDto(newLocation), nil
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
