package service

import (
	"context"
	"errors"
	"ricerise/internal/apperror"
	"ricerise/internal/config"
	"ricerise/internal/dal/query"
	"ricerise/internal/dto/querydto"
	"ricerise/internal/dto/response"
	"ricerise/internal/model"
	"ricerise/internal/repository"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
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

	userCount, _ := a.userRepository.Count(goContext, query.UserModel.ID.Column().Name)
	dinnerCount, _ := a.dinnerRepository.Count(goContext, query.DinnerModel.ID.Column().Name)
	currentDinner, _ := a.dinnerRepository.Where(query.DinnerModel.Status.Eq(model.DINNER_ONGOING)).Count(goContext, query.DinnerModel.ID.Column().Name)
	locationCount, _ := a.dinnerRepository.Count(goContext, query.LocationModel.ID.Column().Name)

	return &response.AdminStatusResponse{
		CurrentDinner: int(currentDinner),
		TotalDinner:   int(dinnerCount),
		TotalLocation: int(locationCount),
		TotalUser:     int(userCount),
	}
}

func (a *AdminService) ReviewComment(ctx *gin.Context, id uint64, pass bool) error {
	goContext := ctx.Request.Context()
	_, err := a.commentRepository.Where(query.CommentModel.ID.Eq(id)).First(goContext)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.AccessNoFoundError
		}
		return err
	}

	_, err = a.commentRepository.Where(query.CommentModel.ID.Gt(id)).Set(query.CommentModel.Reviewed.Set(pass)).Update(goContext)
	if err != nil {
		return err
	}
	return nil
}

func (a *AdminService) ReviewLocation(ctx *gin.Context, id uint64, pass bool) error {
	goContext := ctx.Request.Context()
	_, err := a.locationRepository.Where(query.LocationModel.ID.Eq(id)).First(goContext)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.AccessNoFoundError
		}
		return err
	}
	_, err = a.locationRepository.Where(query.LocationModel.ID.Eq(id)).Set(query.LocationModel.Reviewed.Set(pass)).Update(goContext)
	if err != nil {
		return err
	}
	return nil
}

func getReviewList[T any](ctx *gin.Context, pageQuery querydto.AdminPageQuery, repo interface {
	FindNotReviewedOrderedByCreatedAt(ctx context.Context, pageSize int, startId *uint64, startTime *time.Time) ([]*T, error)
}) ([]*T, bool, error) {
	result, err := repo.FindNotReviewedOrderedByCreatedAt(ctx, pageQuery.PageSize, pageQuery.StartId, pageQuery.StartTime)
	if err != nil {
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

func (a *AdminService) GetLocationReviewList(ctx *gin.Context, pageQuery querydto.AdminPageQuery) ([]*model.LocationModel, bool, error) {
	return getReviewList[model.LocationModel](ctx, pageQuery, a.locationRepository)
}

func (a *AdminService) GetCommentReviewList(ctx *gin.Context, pageQuery querydto.AdminPageQuery) ([]*model.CommentModel, bool, error) {
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
