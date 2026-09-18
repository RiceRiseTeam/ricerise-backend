package service

import (
	"ricerise/internal/config"
	"ricerise/internal/repository"

	"github.com/samber/do/v2"
)

type DinnerService struct{
	appConfig *config.AppConfig
}

func (h DinnerService) List(context *gin.Context) any{
	
}