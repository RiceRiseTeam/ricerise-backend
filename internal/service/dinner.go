package service

import (
	"context"
	"ricerise/internal/config"
	"ricerise/internal/model"
	"ricerise/internal/repository"

	"github.com/samber/do/v2"
)

type DinnerService struct {
	dinnerRepo      *repository.DinnerRepository
	participantRepo *repository.ParticipantRepository
	appConfig       *config.AppConfig
}

func (s DinnerService) List(ctx context.Context) (*[]model.DinnerModel, error) {
	//return s.dinnerRepo.FindAllByA(ctx, &model.DinnerModel{Status: 0}, func(db *gorm.DB) *gorm.DB {
	//	return db.Preload("Location").Preload("Host")
	//})
	return nil, nil
}

func (s DinnerService) LikeFind(ctx context.Context, locationName string) (*[]model.DinnerModel, error) {
	//return s.dinnerRepo.FindByLocationName(ctx, locationName)
	return nil, nil
}

func (s DinnerService) ListParticipants(ctx context.Context, dinnerID uint64) (*[]model.ParticipantModel, error) {
	//return s.participantRepo.ListParticipant(ctx, dinnerID)
	return nil, nil
}

func (s DinnerService) NewParticipate(ctx context.Context, userID, dinnerID uint64) error {
	//if _, err := s.dinnerRepo.FindById(ctx, dinnerID); err != nil {
	//	return err
	//}
	//return s.participantRepo.Create(ctx, &model.ParticipantModel{
	//	UserId:   userID,
	//	DinnerId: dinnerID,
	//	Status:   0,
	//})
	return nil
}

func NewDinnerService(injector do.Injector) (*DinnerService, error) {
	return &DinnerService{
		dinnerRepo:      do.MustInvoke[*repository.DinnerRepository](injector),
		participantRepo: do.MustInvoke[*repository.ParticipantRepository](injector),
		appConfig:       do.MustInvoke[*config.AppConfig](injector),
	}, nil
}
