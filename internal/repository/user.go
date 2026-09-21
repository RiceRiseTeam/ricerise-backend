package repository

import (
	"ricerise/internal/model"

	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type UserRepository struct {
	gorm.Interface[model.UserModel]
	db *gorm.DB
}

func NewUserRepository(injector do.Injector) (*UserRepository, error) {
	db := do.MustInvoke[*gorm.DB](injector)

	return &UserRepository{
		Interface: gorm.G[model.UserModel](db),
		db:        db,
	}, nil
}
