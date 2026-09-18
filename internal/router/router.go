package router

import (
	"ricerise/internal/handler"
	"ricerise/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type Router interface {
	RegisterRouters(router *gin.RouterGroup)
}

func New(injector do.Injector) (*gin.Engine, error) {
	engine := gin.New()

	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	errorMiddleware := do.MustInvoke[*middleware.ErrorMiddleware](injector)
	engine.Use(errorMiddleware.CreateHandler())

	api := engine.Group("/api/v1")

	userHandler := do.MustInvoke[*handler.UserHandler](injector)
	userHandler.RegisterRouters(api)

	adminHandler := do.MustInvoke[*handler.AdminHandler](injector)
	adminHandler.RegisterRouters(api)

	mapHandler := do.MustInvoke[*handler.MapHandler](injector)
	mapHandler.RegisterRouters(api)

	return engine, nil
}
