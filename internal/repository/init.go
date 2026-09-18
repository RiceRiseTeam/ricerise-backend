package repository

import (
	"context"

	"gorm.io/gorm"
)

type QueryArgs func(db *gorm.DB) *gorm.DB

type BaseRepository[T any] struct {
	db *gorm.DB
}

func (r *BaseRepository[T]) Create(context context.Context, entity *T) error {
	return r.db.WithContext(context).Create(entity).Error
}

func (r *BaseRepository[T]) FindByIdA(context context.Context, id int64, queryArgs QueryArgs) (*T, error) {
	var result T
	if err := queryArgs(r.db.WithContext(context)).
		First(&result, id).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *BaseRepository[T]) FindByA(context context.Context, where *T, queryArgs QueryArgs) (*T, error) {
	var result T
	if err := queryArgs(r.db.WithContext(context)).Where(where).First(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *BaseRepository[T]) FindAllByA(context context.Context, where *T, queryArgs QueryArgs) (*[]T, error) {
	var results []T
	if err := queryArgs(r.db.WithContext(context)).Where(where).Find(&results).Error; err != nil {
		return nil, err
	}

	return &results, nil
}

func (r *BaseRepository[T]) FindById(context context.Context, id int64) (*T, error) {
	return r.FindByIdA(context, id, func(db *gorm.DB) *gorm.DB { return db })
}

func (r *BaseRepository[T]) FindBy(context context.Context, where *T) (*T, error) {
	return r.FindByA(context, where, func(db *gorm.DB) *gorm.DB { return db })
}

func (r *BaseRepository[T]) LikeFindBy(context context.Context,keyword string,fields []string,queryArgs QueryArgs,
) (*[]T, error) {
    var results []T
    if keyword == "" {
        if err := queryArgs(r.db.WithContext(ctx)).Find(&results).Error; err != nil {
            return nil, err
        }
        return &results, nil
    }
    pattern := "%" + keyword + "%"
    var conditions []string
    var args []interface{}
    for _, field := range fields {
        conditions = append(conditions, field+" LIKE ?")
        args = append(args, pattern)
    }
    db := r.db.WithContext(ctx).Where(strings.Join(conditions, " OR "), args...)
    if err := queryArgs(db).Find(&results).Error; err != nil {
        return nil, err
    }
    return &results, nil
}