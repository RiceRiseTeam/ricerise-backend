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
