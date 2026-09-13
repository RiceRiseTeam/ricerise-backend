package main

import "github.com/gin-gonic/gin"

func main() {
	gin.SetMode(gin.DebugMode)
	engine := gin.New()
	err := engine.Run()
	if err != nil {
		panic(err)
	}
}
