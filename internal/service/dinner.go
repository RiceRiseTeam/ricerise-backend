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
	"ricerise/internal/logger"
	"ricerise/internal/model"
	"ricerise/internal/repository"
	"strconv"
	"time"
	"uuid"

	sse "github.com/dan-sherwin/go-sse"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

const redisInviteCode = "invite:"
const redisDinnerInviteCode = "dinner:invite:"

// dinnerGenerateCodeScript KEYS: [dinner -> code, code -> dinner] args:[dinnerId, inviteCode, ttl, redisInviteCode]
var dinnerGenerateCodeScript = redis.NewScript(`
if redis.call('EXISTS', KEYS[2]) == 1 then
    return 0
end

local oldCode = redis.call('GET', KEYS[1])
if oldCode then
    redis.call('DEL', ARGV[4] .. oldCode)
end

redis.call('SET', KEYS[1], ARGV[2], 'EX', ARGV[3])
redis.call('SET', KEYS[2], ARGV[1], 'EX', ARGV[3])

return 1
`)

type DinnerService struct {
	dinnerRepo         *repository.DinnerRepository
	participantRepo    *repository.ParticipantRepository
	transactionManager *repository.TransactionManager
	userRepo           *repository.UserRepository
	appConfig          *config.AppConfig
	redisClient        *redis.Client
}

func (d DinnerService) GenerateInviteCode(ctx *gin.Context, userId uint64, dinnerId uint64) (*string, error) {
	goContext := ctx.Request.Context()
	dinner, err := d.dinnerRepo.Where(query.DinnerModel.ID.Eq(dinnerId)).First(goContext)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.AccessNoFoundError
		}
		return nil, err
	}

	if dinner.HostId != userId {
		return nil, apperror.NoPermissionError
	}

	var inviteCode = uuid.New().String()

	keys := []string{
		redisDinnerInviteCode + strconv.FormatUint(dinnerId, 10),
		redisInviteCode + inviteCode,
	}

	result, err := dinnerGenerateCodeScript.Run(goContext, d.redisClient, keys, dinnerId, inviteCode, (5 * time.Minute).Seconds(), redisInviteCode).Int()
	if err != nil {
		return nil, err
	}
	if result == 0 {
		return d.GenerateInviteCode(ctx, userId, dinnerId)
	}
	return &inviteCode, nil
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

	//TODO: 确认操作合法性

	_, err = d.dinnerRepo.Where(query.DinnerModel.ID.Eq(dinnerId)).Set(query.DinnerModel.Status.Set(status)).Update(goContext)
	return err
}

func (d DinnerService) JoinDinner(ctx *gin.Context, userId uint64, dinnerId uint64, code string) (*dto.DinnerDto, error) {
	goContext := ctx.Request.Context()
	codeDinnerId, err := d.redisClient.Get(goContext, redisInviteCode+code).Uint64()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, apperror.InviteCodeNotFoundError
		}
		return nil, err
	}

	if codeDinnerId != dinnerId {
		return nil, apperror.InternalServerError
	}

	dinner, err := d.joinDinner(goContext, userId, dinnerId)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.ParticipantConflictError
		}
		return nil, err
	}

	return dto.NewDinnerDto(dinner), nil
}

func (d DinnerService) joinDinner(ctx context.Context, userId uint64, dinnerId uint64) (*model.DinnerModel, error) {
	dinners, err := d.getCurrentDinners(ctx, userId)
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
			Preload(query.DinnerModel.Host.Name(), nil).Preload(query.DinnerModel.Participants.Name(), nil).First(ctx)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.AccessNoFoundError
			}
			return err
		}
		dinner = &targetDinner
		err = d.participantRepo.WithTx(tx).Create(ctx, &model.ParticipantModel{
			UserId:   userId,
			DinnerId: dinnerId,
		})
		if err != nil {
			return err
		}

		count, err := d.participantRepo.WithTx(tx).Where(query.ParticipantModel.DinnerId.Eq(dinnerId)).Count(ctx, query.ParticipantModel.ID.Column().Name)
		if err != nil {
			return err
		}

		if int(count) >= targetDinner.MaxPeople {
			_, err = d.dinnerRepo.WithTx(tx).Where(query.DinnerModel.ID.Eq(targetDinner.ID)).Set(query.DinnerModel.Status.Set(model.DINNER_FULL)).Update(ctx)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	_ = d.broadcastMessage(ctx, userId, dinnerId, "toast", dto.ToastSSE{Message: fmt.Sprintf("用户%s 加入了饭局")})
	return dinner, nil
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
		MaxPeople:  request.MaxPeople,
		MeetTime:   request.MeetTime,
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

func (d DinnerService) GetDinnerList(ctx *gin.Context, userId uint64) ([]*dto.DinnerDto, error) {
	result, err := d.getCurrentDinners(ctx.Request.Context(), userId)
	if err != nil {
		return nil, err
	}
	return dto.Map(result, func(t model.DinnerModel) *dto.DinnerDto {
		return dto.NewDinnerDto(&t)
	}), nil
}

func (d DinnerService) GetDinnerDetail(ctx *gin.Context, dinnerId uint64) (*dto.DinnerDto, error) {
	dinner, err := d.dinnerRepo.Where(query.DinnerModel.ID.Eq(dinnerId)).
		Preload(query.DinnerModel.Participants.Name(), nil).
		Preload(query.DinnerModel.Host.Name(), nil).
		Preload(query.DinnerModel.Location.Name(), nil).First(ctx.Request.Context())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.AccessNoFoundError
		}
		return nil, err
	}

	return dto.NewDinnerDto(&dinner), nil
}

func (d DinnerService) getCurrentDinners(ctx context.Context, userId uint64) ([]model.DinnerModel, error) {
	sub := d.participantRepo.
		Select(query.ParticipantModel.DinnerId.Column().Name).
		Where(query.ParticipantModel.UserId.Eq(userId))
	dinners, err := d.dinnerRepo.
		Where(query.DinnerModel.Status.Neq(model.DINNER_CANCELLED), query.DinnerModel.Status.Neq(model.DINNER_FINISHED)).
		Where("id IN (?)", sub).Preload(query.DinnerModel.Participants.Name(), nil).Preload(query.DinnerModel.Host.Name(), nil).Preload(query.DinnerModel.Location.Name(), nil).
		Find(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	logger.Info("get dinners", dinners)

	return dinners, nil
}

func NewDinnerService(injector do.Injector) (*DinnerService, error) {
	return &DinnerService{
		dinnerRepo:         do.MustInvoke[*repository.DinnerRepository](injector),
		participantRepo:    do.MustInvoke[*repository.ParticipantRepository](injector),
		transactionManager: do.MustInvoke[*repository.TransactionManager](injector),
		userRepo:           do.MustInvoke[*repository.UserRepository](injector),
		appConfig:          do.MustInvoke[*config.AppConfig](injector),
		redisClient:        do.MustInvoke[*redis.Client](injector),
	}, nil
}
