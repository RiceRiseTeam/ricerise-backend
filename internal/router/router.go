package router

import (
	"ricerise/internal/handler"
	"ricerise/internal/middleware"
	"time"

	_ "ricerise/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router interface {
	RegisterRouters(router *gin.RouterGroup)
}

func New(injector do.Injector) (*gin.Engine, error) {
	engine := gin.New()

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:7891", "https://rice.linstar.cn"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Accept", "Authorization",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	errorMiddleware := do.MustInvoke[*middleware.ErrorMiddleware](injector)
	engine.Use(errorMiddleware.CreateHandler())

	api := engine.Group("/api/v1")

	userHandler := do.MustInvoke[*handler.UserHandler](injector)
	userHandler.RegisterRouters(api)

	DinnerHandler := do.MustInvoke[*handler.DinnerHandler](injector)
	DinnerHandler.RegisterRouters(api)

	adminHandler := do.MustInvoke[*handler.AdminHandler](injector)
	adminHandler.RegisterRouters(api)

	mapHandler := do.MustInvoke[*handler.MapHandler](injector)
	mapHandler.RegisterRouters(api)

	agentHandler := do.MustInvoke[*handler.AgentHandler](injector)
	agentHandler.RegisterRouters(api)

	return engine, nil
}
