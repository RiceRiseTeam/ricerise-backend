package repository

import (
	"context"
	"ricerise/internal/model"
	"time"

	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type CommentRepository struct {
	BaseRepository[model.CommentModel]
	db *gorm.DB
}

func (l CommentRepository) FindNotReviewedOrderedByCreatedAt(ctx context.Context, pageSize int, startId *uint64, startTime *time.Time) ([]*model.CommentModel, error) {
	var result []*model.CommentModel
	var reviewed = false
	query := l.db.WithContext(ctx).Where(&model.CommentModel{
		Reviewed: &reviewed,
	})

	if startId != nil && startTime != nil {
		query = query.Where("(created_at, id) < (?, ?)", startTime, startId)
	}

	query = query.Order("created_at DESC, id DESC").Limit(pageSize + 1).Find(&result)
	return result, query.Error
}

func NewCommentRepository(injector do.Injector) (*CommentRepository, error) {
	db := do.MustInvoke[*gorm.DB](injector)
	return &CommentRepository{
		BaseRepository: BaseRepository[model.CommentModel]{db: db},
		db:             db,
	}, nil
}
