package handler

import (
	"ricerise/internal/config"
	"ricerise/internal/middleware"
	"ricerise/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type AdminHandler struct {
	appConfig      *config.AppConfig
	authMiddleware *middleware.AuthMiddleware
	adminService   *service.AdminService
}

func (a AdminHandler) RegisterRouters(router *gin.RouterGroup) {
	api := router.Group("/admin")
	api.Use(a.authMiddleware.CreateHandler(a.authMiddleware.AdminLevel))
}

func NewAdminHandler(injector do.Injector) (*AdminHandler, error) {
	authMiddleware := do.MustInvoke[*middleware.AuthMiddleware](injector)
	adminService := do.MustInvoke[*service.AdminService](injector)
	appConfig := do.MustInvoke[*config.AppConfig](injector)
	return &AdminHandler{
		appConfig:      appConfig,
		authMiddleware: authMiddleware,
		adminService:   adminService,
	}, nil
}
