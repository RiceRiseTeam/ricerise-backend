package repository

import (
	"context"
	"ricerise/internal/model"
	"time"

	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type LocationRepository struct {
	BaseRepository[model.LocationModel]
	db *gorm.DB
}

func (l LocationRepository) FindNotReviewedOrderedByCreatedAt(ctx context.Context, pageSize int, startId *uint64, startTime *time.Time) ([]*model.LocationModel, error) {
	var result []*model.LocationModel
	var reviewed = false
	query := l.db.WithContext(ctx).Where(&model.LocationModel{
		Reviewed: &reviewed,
	})

	if startId != nil && startTime != nil {
		query = query.Where("(created_at, id) < (?, ?)", startTime, startId)
	}

	query = query.Order("created_at DESC, id DESC").Limit(pageSize + 1).Find(&result)
	return result, query.Error
}

func NewLocationRepository(injector do.Injector) (*LocationRepository, error) {
	db := do.MustInvoke[*gorm.DB](injector)

	return &LocationRepository{
		BaseRepository: BaseRepository[model.LocationModel]{db: db},
		db:             db,
	}, nil
}
