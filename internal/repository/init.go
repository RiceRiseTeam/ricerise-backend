package repository

import (
	"context"
	"strings"

	"gorm.io/gorm"
)

type QueryArgs func(db *gorm.DB) *gorm.DB

type BaseRepository[T any] struct {
	db *gorm.DB
}

func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *BaseRepository[T]) FindByIdA(ctx context.Context, id int64, queryArgs QueryArgs) (*T, error) {
	var result T
	if err := queryArgs(r.db.WithContext(ctx)).
		First(&result, id).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *BaseRepository[T]) FindByA(ctx context.Context, where *T, queryArgs QueryArgs) (*T, error) {
	var result T
	if err := queryArgs(r.db.WithContext(ctx)).Where(where).First(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *BaseRepository[T]) FindAllByA(ctx context.Context, where *T, queryArgs QueryArgs) (*[]T, error) {
	var results []T
	if err := queryArgs(r.db.WithContext(ctx)).Where(where).Find(&results).Error; err != nil {
		return nil, err
	}

	return &results, nil
}

func (r *BaseRepository[T]) FindById(ctx context.Context, id int64) (*T, error) {
	return r.FindByIdA(ctx, id, func(db *gorm.DB) *gorm.DB { return db })
}

func (r *BaseRepository[T]) FindBy(ctx context.Context, where *T) (*T, error) {
	return r.FindByA(ctx, where, func(db *gorm.DB) *gorm.DB { return db })
}

func (r *BaseRepository[T]) CountBy(ctx context.Context, where *T) int64 {
	var count int64
	r.db.WithContext(ctx).Where(where).Count(&count)
	return count
}

func (r *BaseRepository[T]) Count(ctx context.Context) int64 {
	return r.CountBy(ctx, new(T))
}

func (r *BaseRepository[T]) LikeFindBy(ctx context.Context, keyword string, fields []string, queryArgs QueryArgs,
) (*[]T, error) {
	var results []T
	db := queryArgs(r.db.WithContext(ctx))
	if keyword != "" && len(fields) > 0 {
		pattern := "%" + keyword + "%"
		conditions := make([]string, 0, len(fields))
		args := make([]any, 0, len(fields))
		for _, field := range fields {
			conditions = append(conditions, field+" LIKE ?")
			args = append(args, pattern)
		}
		db = db.Where(strings.Join(conditions, " OR "), args...)
	}
	if err := db.Find(&results).Error; err != nil {
		return nil, err
	}
	return &results, nil
}
