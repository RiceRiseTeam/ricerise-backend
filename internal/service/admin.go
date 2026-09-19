package service

import (
	"ricerise/internal/config"
	"ricerise/internal/dto/query"
	"ricerise/internal/dto/response"
	"ricerise/internal/model"
	"ricerise/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type AdminService struct {
	userRepository     *repository.UserRepository
	dinnerRepository   *repository.DinnerRepository
	locationRepository *repository.LocationRepository
	appConfig          *config.AppConfig
}

func (a *AdminService) GetAppStatus(ctx *gin.Context) *response.AdminStatusResponse {
	goContext := ctx.Request.Context()

	userCount := a.userRepository.Count(goContext)
	dinnerCount := a.dinnerRepository.Count(goContext)
	currentDinner := a.dinnerRepository.CountBy(goContext, &model.DinnerModel{
		Status: model.DINNER_ONGOING,
	})
	locationCount := a.dinnerRepository.Count(goContext)

	return &response.AdminStatusResponse{
		CurrentDinner: int(currentDinner),
		TotalDinner:   int(dinnerCount),
		TotalLocation: int(locationCount),
		TotalUser:     int(userCount),
	}
}

func (a *AdminService) GetLocationsReviewList(ctx *gin.Context, pageQuery query.AdminPageQuery) ([]*model.LocationModel, bool, error) {
	result, err := a.locationRepository.FindNotReviewedOrderedByCreatedAt(ctx, pageQuery.PageSize, pageQuery.StartId, pageQuery.StartTime)
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

func (a *AdminService) GetCommentReviewList(ctx *gin.Context) {

}

func NewAdminService(injector do.Injector) (*AdminService, error) {
	userRepository := do.MustInvoke[*repository.UserRepository](injector)
	dinnerRepository := do.MustInvoke[*repository.DinnerRepository](injector)
	appConfig := do.MustInvoke[*config.AppConfig](injector)
	return &AdminService{
		userRepository:   userRepository,
		dinnerRepository: dinnerRepository,
		appConfig:        appConfig,
	}, nil
}
