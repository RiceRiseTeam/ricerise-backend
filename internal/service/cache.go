package service

import (
	"github.com/redis/go-redis/v9"
	"github.com/samber/do/v2"
)

type CacheService struct {
	redisClient *redis.Client
}

func NewCacheService(injector do.Injector) (*CacheService, error) {
	redisClient := do.MustInvoke[*redis.Client](injector)
	return &CacheService{redisClient: redisClient}, nil
}
