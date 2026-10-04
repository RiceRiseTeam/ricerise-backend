package service

import (
	"context"
	"errors"
	"fmt"
	"ricerise/internal/apperror"
	"ricerise/internal/config"
	"ricerise/internal/dal/query"
	"ricerise/internal/dto"
	"ricerise/internal/dto/request"
	"ricerise/internal/model"
	"ricerise/internal/repository"

	sse "github.com/dan-sherwin/go-sse"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type DinnerService struct {
	dinnerRepo         *repository.DinnerRepository
	participantRepo    *repository.ParticipantRepository
	transactionManager *repository.TransactionManager
	userRepo           *repository.UserRepository
	appConfig          *config.AppConfig
}

func (d DinnerService) UpdateDinnerStatus(ctx *gin.Context, userId uint64, dinnerId uint64, status int8) error {
	goContext := ctx.Request.Context()
	dinner, err := d.dinnerRepo.Where(query.DinnerModel.ID.Eq(dinnerId)).First(goContext)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.AccessNoFoundError
		}
		return err
	}

	if dinner.HostId != userId {
		return apperror.NoPermissionError
	}

	_, err = d.dinnerRepo.Where(query.DinnerModel.ID.Eq(dinnerId)).Set(query.DinnerModel.Status.Set(status)).Update(goContext)
	return err
}

func (d DinnerService) JoinDinner(ctx *gin.Context, userId uint64, dinnerId uint64) (*dto.DinnerDto, error) {
	goContext := ctx.Request.Context()
	dinners, err := d.getCurrentDinners(goContext, userId)
	if err != nil {
		return nil, err
	}
	if len(dinners) > 3 {
		return nil, apperror.DinnerConflictError
	}

	var dinner *model.DinnerModel = nil

	err = d.transactionManager.Do(func(tx *gorm.DB) error {
		targetDinner, err := d.dinnerRepo.WithTx(tx).
			Where(query.DinnerModel.ID.Eq(dinnerId), query.DinnerModel.Status.Eq(model.DINNER_HIRING)).
			Preload(query.DinnerModel.Host.Name(), nil).Preload(query.DinnerModel.Participants.Name(), nil).First(goContext)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.AccessNoFoundError
			}
			return err
		}
		dinner = &targetDinner
		err = d.participantRepo.WithTx(tx).Create(goContext, &model.ParticipantModel{
			UserId:   userId,
			DinnerId: dinnerId,
		})
		if err != nil {
			return err
		}

		count, err := d.participantRepo.WithTx(tx).Where(query.ParticipantModel.DinnerId.Eq(dinnerId)).Count(goContext, query.ParticipantModel.ID.Column().Name)
		if err != nil {
			return err
		}

		if int(count) >= targetDinner.MaxPeople {
			_, err = d.dinnerRepo.WithTx(tx).Where(query.DinnerModel.ID.Eq(targetDinner.ID)).Set(query.DinnerModel.Status.Set(model.DINNER_FULL)).Update(goContext)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	_ = d.broadcastMessage(ctx, userId, dinnerId, "toast", dto.ToastSSE{Message: fmt.Sprintf("用户%s 加入了饭局")})
	return dto.NewDinnerDto(dinner), nil
}

func (d DinnerService) ExistDinner(ctx *gin.Context, userId uint64, dinnerId uint64) error {
	goContext := ctx.Request.Context()
	participated, err := d.participantRepo.IsParticipant(goContext, userId, dinnerId)
	if err != nil {
		return err
	}
	if !participated {
		return apperror.AccessNoFoundError
	}

	user, err := d.userRepo.Where(query.UserModel.ID.Eq(userId)).First(goContext)
	if err != nil {
		return err
	}

	err = d.transactionManager.Do(func(tx *gorm.DB) error {
		dinner, err := d.dinnerRepo.Where(query.DinnerModel.ID.Eq(dinnerId)).First(goContext)
		if err != nil {
			return err
		}

		_, err = d.participantRepo.WithTx(tx).
			Where(query.ParticipantModel.UserId.Eq(userId)).
			Where(query.ParticipantModel.DinnerId.Eq(dinnerId)).Delete(goContext)

		count, err := d.participantRepo.WithTx(tx).Where(query.ParticipantModel.DinnerId.Eq(dinnerId)).Count(goContext, query.ParticipantModel.ID.Column().Name)
		if err != nil {
			return err
		}

		if int(count) < dinner.MaxPeople {
			_, err = d.dinnerRepo.WithTx(tx).
				Where(query.DinnerModel.ID.Eq(dinnerId), query.DinnerModel.Status.Eq(model.DINNER_FULL)).
				Set(query.DinnerModel.Status.Set(model.DINNER_HIRING)).Update(goContext)
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	_ = d.broadcastMessage(goContext, userId, dinnerId, "toast", dto.ToastSSE{Message: fmt.Sprintf("用户%s 退出了饭局", user.Nickname)})

	return nil
}

func (d DinnerService) SendDinnerMessage(ctx *gin.Context, userId uint64, dinnerId uint64, message string) error {
	goContext := ctx.Request.Context()
	user, err := d.userRepo.Where(query.UserModel.ID.Eq(userId)).First(goContext)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NoAccessTokenError
		}
		return err
	}
	err = d.broadcastMessage(ctx.Request.Context(), userId, dinnerId, "chat", &dto.ChatSSE{User: dto.NewUserDto(&user), Message: message})
	if err != nil {
		return err
	}

	return nil
}

func (d DinnerService) broadcastMessage(ctx context.Context, userId uint64, dinnerId uint64, eventType string, data any) error {
	participated, err := d.participantRepo.IsParticipant(ctx, userId, dinnerId)
	if err != nil {
		return err
	}
	if !participated {
		return apperror.AccessNoFoundError
	}

	result, err := d.participantRepo.
		Preload(query.ParticipantModel.User.Name(), nil).
		Where(query.ParticipantModel.DinnerId.Eq(dinnerId)).
		Find(ctx)
	if err != nil {
		return err
	}

	for _, participant := range result {
		_ = sse.SendEventToUID(eventType, data, participant.User.Username)
	}
	return nil
}

func (d DinnerService) CreateDinner(ctx *gin.Context, userId uint64, request request.CreateDinnerRequest) (*dto.DinnerDto, error) {
	goContext := ctx.Request.Context()
	dinner, err := d.getCurrentDinners(goContext, userId)
	if err != nil {
		return nil, err
	}
	if len(dinner) > 3 {
		return nil, apperror.DinnerConflictError
	}
	newDinner := &model.DinnerModel{
		LocationId: request.LocationId,
		HostId:     userId,
	}

	err = d.dinnerRepo.Create(goContext, newDinner)
	if err != nil {
		return nil, err
	}

	newParticipant := &model.ParticipantModel{
		UserId:   userId,
		DinnerId: newDinner.ID,
	}
	err = d.participantRepo.Create(goContext, newParticipant)

	result, err := d.dinnerRepo.
		Preload(query.DinnerModel.Host.Name(), nil).
		Preload(query.DinnerModel.Participants.Name(), nil).
		Where(query.DinnerModel.ID.Eq(newDinner.ID)).First(goContext)

	if err != nil {
		return nil, err
	}

	return dto.NewDinnerDto(&result), nil
}

func (d DinnerService) getCurrentDinners(ctx context.Context, userId uint64) ([]model.DinnerModel, error) {
	sub := d.participantRepo.
		Select(query.ParticipantModel.DinnerId.Column().Name).
		Where(query.ParticipantModel.UserId.Eq(userId))
	dinners, err := d.dinnerRepo.
		Where(query.DinnerModel.Status.Neq(model.DINNER_CANCELLED), query.DinnerModel.Status.Neq(model.DINNER_FINISHED)).
		Where("id IN (?)", sub).
		Find(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return dinners, nil
}

func NewDinnerService(injector do.Injector) (*DinnerService, error) {
	return &DinnerService{
		dinnerRepo:      do.MustInvoke[*repository.DinnerRepository](injector),
		participantRepo: do.MustInvoke[*repository.ParticipantRepository](injector),
		userRepo:        do.MustInvoke[*repository.UserRepository](injector),
		appConfig:       do.MustInvoke[*config.AppConfig](injector),
	}, nil
}
