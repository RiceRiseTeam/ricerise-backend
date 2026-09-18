package repository

import (
	"ricerise/internal/model"

	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type DinnerRepository struct {
	BaseRepository[model.DinnerModel]
	db *gorm.DB
}

func NewDinnerRepository(injector do.Injector) (*DinnerRepository, error) {
	db := do.MustInvoke[*gorm.DB](injector)

	return &DinnerRepository{
		BaseRepository: BaseRepository[model.DinnerModel]{db: db},
		db:             db,
	}, nil
}

func (h DinnerRepository)ListParticipant(context context.Context, id int64) (*DinnerModel, error) {
	return h.FindByIDA(context, id, func(db *gorm.DB) *gorm.DB {
		return db.Select("User")}).Preload("UserId")
}