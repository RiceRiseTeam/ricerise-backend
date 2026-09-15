package main

import (
	"fmt"
	"ricerise/internal/config"
	"ricerise/internal/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)
	engine := gin.New()
	logger.Info("starting...")
	cfg := config.Get()

	err := engine.Run(fmt.Sprintf("127.0.0.1:%d", cfg.AppPort))
	if err != nil {
		panic(err)
	}
}
