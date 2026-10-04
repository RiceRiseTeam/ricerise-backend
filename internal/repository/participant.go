package repository

import (
	"context"
	"errors"
	"ricerise/internal/dal/query"
	"ricerise/internal/model"

	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type ParticipantRepository struct {
	gorm.Interface[model.ParticipantModel]
	db *gorm.DB
}

func (p ParticipantRepository) WithTx(tx *gorm.DB) *ParticipantRepository {
	return &ParticipantRepository{
		Interface: gorm.G[model.ParticipantModel](tx),
		db:        tx,
	}
}

func (p ParticipantRepository) IsParticipant(ctx context.Context, userId uint64, dinnerId uint64) (bool, error) {
	_, err := p.Interface.Where(query.ParticipantModel.DinnerId.Eq(dinnerId), query.ParticipantModel.UserId.Eq(userId)).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func NewParticipantRepository(injector do.Injector) (*ParticipantRepository, error) {
	db := do.MustInvoke[*gorm.DB](injector)

	return &ParticipantRepository{
		Interface: gorm.G[model.ParticipantModel](db),
		db:        db,
	}, nil
}

//func (r ParticipantRepository) ListParticipant(ctx context.Context, dinnerID uint64) (*[]model.ParticipantModel, error) {
//	return r.FindAllByA(ctx, &model.ParticipantModel{
//		DinnerId: dinnerID,
//	}, func(db *gorm.DB) *gorm.DB {
//		return db.Select("id", "dinner_id", "user_id", "status").Preload("User")
//	})
//}
