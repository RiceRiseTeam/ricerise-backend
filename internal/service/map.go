package service

import (
	"context"
	"errors"
	"ricerise/internal/apperror"
	"ricerise/internal/config"
	"ricerise/internal/dal/query"
	"ricerise/internal/dto"
	"ricerise/internal/dto/request"
	"ricerise/internal/logger"
	"ricerise/internal/middleware"
	"ricerise/internal/model"
	"ricerise/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/restayway/gogis"
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
	goContext := ctx.Request.Context()
	comment, err := m.fetchComment(goContext, id)
	if err != nil {
		return err
	}

	userInfo := m.authMiddleware.GetUserInfo(ctx)
	if userInfo == nil {
		return apperror.NoAccessTokenError
	}
	if userInfo.PermissionLevel < m.authMiddleware.OwnerLevel && userInfo.Username == comment.User.Username {
		return apperror.NoPermissionError
	}
	_, err = m.commentRepository.Where(query.CommentModel.ID.Eq(id)).Delete(goContext)
	if err != nil {
		return err
	}
	return nil
}

func (m MapService) DeleteLocation(ctx *gin.Context, id uint64) error {
	goContext := ctx.Request.Context()
	_, err := m.fetchLocation(goContext, id)
	if err != nil {
		return err
	}
	userInfo := m.authMiddleware.GetUserInfo(ctx)
	if userInfo == nil {
		return apperror.NoAccessTokenError
	}
	if userInfo.PermissionLevel < m.authMiddleware.OwnerLevel {
		return apperror.NoPermissionError
	}

	_, err = m.locationRepository.Where(query.LocationModel.ID.Eq(id)).Delete(goContext)
	if err != nil {
		return err
	}
	return nil
}

func (m MapService) GetLocationDetail(ctx *gin.Context, id uint64) (*dto.LocationDto, error) {
	location, err := m.locationRepository.Where(query.LocationModel.ID.Eq(id)).Preload(query.LocationModel.User.Name(), nil).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.AccessNoFoundError
		}
		return nil, err
	}
	return dto.NewLocationDto(&location), nil
}

func (m MapService) GetLocationsInRange(ctx *gin.Context, request request.GetLocationsRequest) ([]*dto.LocationDto, error) {
	result, err := m.locationRepository.FindAllInRange(ctx.Request.Context(), request.MinLng, request.MinLat, request.MaxLng, request.MaxLat)
	if err != nil {
		return nil, err
	}
	logger.Info("getLocationsInRange", len(result))

	return dto.Map(result, func(t model.LocationModel) *dto.LocationDto {
		return dto.NewLocationDto(&t)
	}), nil
}

func (m MapService) fetchLocation(ctx context.Context, id uint64) (*model.LocationModel, error) {
	location, err := m.locationRepository.Where(query.LocationModel.ID.Eq(id)).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.AccessNoFoundError
		}
		return nil, err
	}
	return &location, nil
}

func (m MapService) fetchComment(ctx context.Context, id uint64) (*model.CommentModel, error) {
	comment, err := m.commentRepository.Where(query.LocationModel.ID.Eq(id)).Preload(query.CommentModel.User.Name(), nil).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.AccessNoFoundError
		}
		return nil, err
	}
	return &comment, nil
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
		Location: gogis.Point{
			Lng: request.Longitude,
			Lat: request.Latitude,
		},
		Name:        request.Name,
		Description: request.Description,
		UserId:      userInfo.UserId,
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
