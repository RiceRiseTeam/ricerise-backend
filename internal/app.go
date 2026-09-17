package internal

import (
	"fmt"
	"ricerise/internal/cache"
	"ricerise/internal/config"
	"ricerise/internal/database"
	"ricerise/internal/handler"
	"ricerise/internal/logger"
	"ricerise/internal/middleware"
	"ricerise/internal/repository"
	"ricerise/internal/router"
	"ricerise/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

var engine *gin.Engine
var cfg *config.AppConfig

func init() {
	gin.SetMode(gin.DebugMode)
	logger.Info("starting...")

	root := do.New()
	// 基础设施
	do.Provide(root, config.New)
	do.Provide(root, database.NewMySQL)
	do.Provide(root, cache.NewRedis)
	// Repo 层
	do.Provide(root, repository.NewUserRepository)

	// Service 层
	do.Provide(root, service.NewUserService)
	do.Provide(root, service.NewCacheService)

	// Handler 层
	do.Provide(root, handler.NewUserHandler)

	// 中间件
	do.Provide(root, middleware.NewErrorMiddleware)
	do.Provide(root, middleware.NewAuthMiddleware)

	// 路由注册
	do.Provide(root, router.New)

	engine = do.MustInvoke[*gin.Engine](root)
	cfg = do.MustInvoke[*config.AppConfig](root)
}

func Start() {
	err := engine.Run(fmt.Sprintf("127.0.0.1:%d", cfg.AppPort))
	if err != nil {
		panic(err)
	}
}

func GetEngine() *gin.Engine {
	return engine
}
