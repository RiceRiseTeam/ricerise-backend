package repository

import (
	"context"
	"ricerise/internal/model"

	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type ParticipantRepository struct {
	BaseRepository[model.ParticipantModel]
	db *gorm.DB
}

func NewParticipantRepository(injector do.Injector) (*ParticipantRepository, error) {
	db := do.MustInvoke[*gorm.DB](injector)

	return &ParticipantRepository{
		BaseRepository: BaseRepository[model.ParticipantModel]{db: db},
		db:             db,
	}, nil
}

func (r ParticipantRepository) ListParticipant(ctx context.Context, dinnerID uint64) (*[]model.ParticipantModel, error) {
	return r.FindAllByA(ctx, &model.ParticipantModel{
		DinnerId: dinnerID,
	}, func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "dinner_id", "user_id", "status").Preload("User")
	})
}
