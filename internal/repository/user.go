package repository

import (
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(injector do.Injector) (*UserRepository, error) {
	db := do.MustInvoke[*gorm.DB](injector)
	return &UserRepository{db: db}, nil
}
