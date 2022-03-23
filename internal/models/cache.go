package models

import (
	"kms/internal/models/redis"
	"time"
)

// CacheType exports
type CacheType = redis.CacheType

// Cache cacheOper
type Cache interface {
	Cache(ct CacheType, key string, cached interface{}, duration ...time.Duration) error
	Query(ct CacheType, key string, entity interface{}, duration ...time.Duration) error
	Del(ct CacheType, key ...string) error
}
