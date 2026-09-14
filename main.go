package main

import (
	"ricerise/internal/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)
	engine := gin.New()
	logger.Info("starting...")
	err := engine.Run()
	if err != nil {
		panic(err)
	}
}
