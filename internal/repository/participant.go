package repository

import (
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
