package repository

import (
	"ricerise/internal/model"

	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type LocationRepository struct {
	BaseRepository[model.LocationModel]
	db *gorm.DB
}

func NewLocationRepository(injector do.Injector) (*LocationRepository, error) {
	db := do.MustInvoke[*gorm.DB](injector)

	return &LocationRepository{
		BaseRepository: BaseRepository[model.LocationModel]{db: db},
		db:             db,
	}, nil
}
