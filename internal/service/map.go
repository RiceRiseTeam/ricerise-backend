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
	comment := m.fetchLocation(ctx, id)
	if comment == nil {
		return apperror.PlaceHolder
	}

	//TODO(linstarowo): 鉴权
	err := m.commentRepository.Delete(ctx, id)
	if err != nil {
		_ = ctx.Error(err)
		return err
	}
	return nil
}

func (m MapService) DeleteLocation(ctx *gin.Context, id uint64) error {
	location := m.fetchLocation(ctx, id)
	if location == nil {
		return apperror.PlaceHolder
	}
	//TODO(linstarowo): 鉴权
	err := m.locationRepository.Delete(ctx, id)
	if err != nil {
		_ = ctx.Error(err)
		return err
	}
	return nil
}

// 此方法会进行错误处理 只需向handler 回传占位符
func (m MapService) fetchLocation(ctx *gin.Context, id uint64) *model.LocationModel {
	location, err := m.locationRepository.FindById(ctx.Request.Context(), id)
	if err != nil {
		_ = ctx.Error(err)
		return nil
	}
	if location == nil {
		_ = ctx.Error(apperror.AccessNoFoundError)
		return nil
	}
	return location
}

// 此方法会进行错误处理 只需向handler 回传占位符
func (m MapService) fetchComment(ctx *gin.Context, id uint64) *model.CommentModel {
	comment, err := m.commentRepository.FindById(ctx.Request.Context(), id)
	if err != nil {
		_ = ctx.Error(err)
		return nil
	}
	if comment == nil {
		_ = ctx.Error(apperror.AccessNoFoundError)
		return nil
	}
	return comment
}

func (m MapService) UploadComment(ctx *gin.Context, request request.UploadCommentRequest) *dto.CommentDto {
	userInfo := m.authMiddleware.GetUserInfo(ctx)
	if userInfo == nil {
		_ = ctx.Error(apperror.NoAccessTokenError)
		return nil
	}
	goContext := ctx.Request.Context()
	location := m.fetchLocation(ctx, request.LocationId)
	if location == nil {
		return nil
	}

	newComment := &model.CommentModel{
		Rating:     request.Rating,
		Content:    request.Content,
		UserId:     userInfo.UserId,
		LocationId: location.ID,
	}
	err := m.commentRepository.Create(goContext, newComment)
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
