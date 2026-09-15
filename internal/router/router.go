package router

import (
	"ricerise/internal/handler"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type Router interface {
	RegisterRouters(router *gin.RouterGroup)
}

func New(injector do.Injector) (*gin.Engine, error) {
	engine := gin.Default()
	api := engine.Group("/api")
	userHandler := do.MustInvoke[*handler.UserHandler](injector)
	userHandler.RegisterRouters(api)
	return engine, nil
}
