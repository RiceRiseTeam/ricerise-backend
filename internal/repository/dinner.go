package repository

import (
	"context"
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

func (r DinnerRepository) FindByLocationName(ctx context.Context, name string) (*[]model.DinnerModel, error) {
	// TODO(stywlkj): Need to be modified
	return nil, nil
	//var results []model.DinnerModel
	//result, err := r.LikeFindBy(ctx, name, []string"neme", results)
	//if err != nil {
	//	return nil ,err
	//}
	//return result, nil
}
