package service

import (
	"context"
	"ricerise/internal/apperror"
	"ricerise/internal/config"
	"ricerise/internal/dto/query"
	"ricerise/internal/dto/response"
	"ricerise/internal/model"
	"ricerise/internal/repository"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type AdminService struct {
	userRepository     *repository.UserRepository
	dinnerRepository   *repository.DinnerRepository
	locationRepository *repository.LocationRepository
	commentRepository  *repository.CommentRepository
	appConfig          *config.AppConfig
}

func (a *AdminService) GetAppStatus(ctx *gin.Context) *response.AdminStatusResponse {
	goContext := ctx.Request.Context()

	userCount, _ := a.userRepository.Count(goContext)
	dinnerCount, _ := a.dinnerRepository.Count(goContext)
	currentDinner, _ := a.dinnerRepository.CountBy(goContext, &model.DinnerModel{
		Status: model.DINNER_ONGOING,
	})
	locationCount, _ := a.dinnerRepository.Count(goContext)

	return &response.AdminStatusResponse{
		CurrentDinner: int(currentDinner),
		TotalDinner:   int(dinnerCount),
		TotalLocation: int(locationCount),
		TotalUser:     int(userCount),
	}
}

func (a *AdminService) ReviewComment(ctx *gin.Context, id uint64, pass bool) error {
	_, err := a.commentRepository.FindById(ctx, id)
	if err != nil {
		_ = ctx.Error(apperror.AccessNoFoundError)
		return err
	}
	err = a.commentRepository.Updates(ctx, &model.CommentModel{ID: id}, &model.CommentModel{Reviewed: &pass})
	if err != nil {
		_ = ctx.Error(err)
		return err
	}
	return nil
}

func (a *AdminService) ReviewLocation(ctx *gin.Context, id uint64, pass bool) error {
	_, err := a.locationRepository.FindById(ctx, id)
	if err != nil {
		_ = ctx.Error(apperror.AccessNoFoundError)
		return err
	}
	err = a.locationRepository.Updates(ctx, &model.LocationModel{ID: id}, &model.LocationModel{Reviewed: &pass})
	if err != nil {
		_ = ctx.Error(err)
		return err
	}
	return nil
}

func getReviewList[T any](ctx *gin.Context, pageQuery query.AdminPageQuery, repo interface {
	FindNotReviewedOrderedByCreatedAt(ctx context.Context, pageSize int, startId *uint64, startTime *time.Time) ([]*T, error)
}) ([]*T, bool, error) {
	result, err := repo.FindNotReviewedOrderedByCreatedAt(ctx, pageQuery.PageSize, pageQuery.StartId, pageQuery.StartTime)
	if err != nil {
		_ = ctx.Error(err)
		return nil, false, err
	}

	var hasNextPage = true
	if len(result) < pageQuery.PageSize {
		hasNextPage = false
	} else {
		result = result[:pageQuery.PageSize]
	}
	return result, hasNextPage, nil
}

func (a *AdminService) GetLocationReviewList(ctx *gin.Context, pageQuery query.AdminPageQuery) ([]*model.LocationModel, bool, error) {
	return getReviewList[model.LocationModel](ctx, pageQuery, a.locationRepository)
}

func (a *AdminService) GetCommentReviewList(ctx *gin.Context, pageQuery query.AdminPageQuery) ([]*model.CommentModel, bool, error) {
	return getReviewList[model.CommentModel](ctx, pageQuery, a.commentRepository)
}

func NewAdminService(injector do.Injector) (*AdminService, error) {
	userRepository := do.MustInvoke[*repository.UserRepository](injector)
	dinnerRepository := do.MustInvoke[*repository.DinnerRepository](injector)
	commentRepository := do.MustInvoke[*repository.CommentRepository](injector)
	locationRepository := do.MustInvoke[*repository.LocationRepository](injector)
	appConfig := do.MustInvoke[*config.AppConfig](injector)
	return &AdminService{
		userRepository:     userRepository,
		dinnerRepository:   dinnerRepository,
		locationRepository: locationRepository,
		commentRepository:  commentRepository,
		appConfig:          appConfig,
	}, nil
}
