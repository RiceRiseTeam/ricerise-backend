package cache

import (
	"context"
	"ricerise/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/samber/do/v2"
)

func NewRedis(injector do.Injector) (*redis.Client, error) {
	appConfig := do.MustInvoke[*config.AppConfig](injector)
	redisClient := redis.NewClient(&redis.Options{
		Addr:         appConfig.RedisUrl,
		Username:     appConfig.RedisUser,
		Password:     appConfig.RedisPassword,
		DB:           1,
		PoolSize:     20,
		MinIdleConns: 5,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		_ = redisClient.Close()
		panic("fail to connect redis: " + err.Error())
	}
	return redisClient, nil
}
