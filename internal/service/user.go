package service

import (
	"ricerise/internal/dto/request"
	"ricerise/internal/repository"

	"github.com/samber/do/v2"
)

type UserService struct {
	repo *repository.UserRepository
}

func (h UserService) RegisterNew(request *request.UserRegisterRequest) {

}

func NewUserService(injector do.Injector) (*UserService, error) {
	repo := do.MustInvoke[*repository.UserRepository](injector)
	return &UserService{repo: repo}, nil
}
