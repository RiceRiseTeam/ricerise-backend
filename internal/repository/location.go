package repository

import (
	"context"
	"ricerise/internal/model"
	"time"

	"github.com/samber/do/v2"
	"golang.org/x/crypto/openpgp/errors"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

type LocationRepository struct {
	gorm.Interface[model.LocationModel]
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

func (l LocationRepository) FindAllInRange(ctx context.Context, minLng, minLat, maxLng, maxLat float64) ([]model.LocationModel, error) {
	if minLng >= maxLng || minLat >= maxLat {
		return nil, errors.InvalidArgumentError("invalid location range")
	}
	return l.Interface.Where(field.NewUnsafeFieldRaw("location @ ST_MakeBox2D(ST_MakePoint(?, ?), ST_MakePoint(?, ?))", minLng, minLat, maxLng, maxLat)).Find(ctx)
}

func NewLocationRepository(injector do.Injector) (*LocationRepository, error) {
	db := do.MustInvoke[*gorm.DB](injector)

	return &LocationRepository{
		Interface: gorm.G[model.LocationModel](db),
		db:        db,
	}, nil
}
