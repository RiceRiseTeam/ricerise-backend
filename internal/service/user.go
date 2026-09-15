package service

import (
	"ricerise/internal/repository"

	"github.com/samber/do/v2"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(injector do.Injector) (*UserService, error) {
	repo := do.MustInvoke[*repository.UserRepository](injector)
	return &UserService{repo: repo}, nil
}
