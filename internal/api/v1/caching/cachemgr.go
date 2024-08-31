package caching

import (
	"github.com/LittleAksMax/blog-backend/internal/cache"
	"github.com/redis/go-redis/v9"
)

type CacheManager struct {
	rdb *redis.Client
}

func NewCacheManager(cacheCfg *cache.Config) *CacheManager {
	return &CacheManager{
		rdb: cacheCfg.Client,
	}
}
