package cache

import (
	"github.com/redis/go-redis/v9"
	"strconv"
)

type RedisCacheService struct {
	Rdb   redis.UniversalClient
	Cache *RedisUtil
}

func NewRedisCacheService(config *RedisConfig) CacheInterface {
	redisUtils := NewRedisUtil(config)
	return &RedisCacheService{Rdb: redisUtils.GetRedisClient(), Cache: redisUtils}
}

func NewRedisCacheServiceByRDB(redisClient redis.UniversalClient) CacheInterface {
	redisUtils := NewRedisUtilByExistRDB(redisClient)
	return &RedisCacheService{Cache: redisUtils}
}

func (l *RedisCacheService) Get(key string) string {
	return l.Cache.Get(key)
}

func (l *RedisCacheService) Set(key string, val string, expiresInSeconds int) {
	l.Cache.Set(key, val, expiresInSeconds)
}

func (l *RedisCacheService) Delete(key string) {
	l.Cache.Delete(key)
}

func (l *RedisCacheService) Exists(key string) bool {
	return l.Cache.Exists(key)
}

func (l *RedisCacheService) GetType() string {
	return "redis"
}

func (l *RedisCacheService) Increment(key string, val int) int {
	cacheVal := l.Cache.Get(key)
	num, err := strconv.Atoi(cacheVal)
	if err != nil {
		num = 0
	}

	ret := num + val

	l.Cache.Set(key, strconv.Itoa(ret), 0)
	return ret
}
