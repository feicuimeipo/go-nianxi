package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coocood/freecache"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type RedisCacheService struct {
	config  *RedisConfig
	redisdb redis.UniversalClient
}

func (r *RedisCacheService) CacheType() string {
	return CACHE_MODE_REDIS
}

func (r *RedisCacheService) GetLocalCache() *freecache.Cache {
	return nil
}

func (r *RedisCacheService) GetRedisClient() redis.UniversalClient {
	return r.redisdb
}

func NewRedis(config *RedisConfig, logger *zap.Logger) ICache {
	rdb := NewRedisClient(config, logger)
	return &RedisCacheService{redisdb: rdb}
}

func (r *RedisCacheService) Set(key string, val interface{}, expire time.Duration) error {
	json, err := json.Marshal(val)
	if err != nil {
		return err
	} else {
		r.redisdb.Set(context.Background(), key, json, expire)
	}

	return nil
}

func (r *RedisCacheService) Expire(key string, expire time.Duration) bool {
	r.redisdb.Expire(context.Background(), key, expire)
	return true
}

func (r *RedisCacheService) SetDefault(key string, value interface{}) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	} else {
		r.redisdb.Set(context.Background(), key, string(val), DefaultExpiration)
	}
	return nil
}

func (r *RedisCacheService) GetDefault(key string, defaultValue interface{}) interface{} {
	if val, exist := r.Get(key); exist {
		return val
	}
	return defaultValue
}

func (r *RedisCacheService) Get(key string) (interface{}, bool) {
	strCmd := r.redisdb.Get(context.Background(), key)
	val, err := strCmd.Result()
	if err == redis.Nil {
		return nil, false
	} else if err != nil {
		return nil, false
	}

	var ret interface{}
	err = json.Unmarshal([]byte(val), &ret)
	if err != nil {
		fmt.Sprintf("%v\n", err)
		return nil, false
	}
	return ret, true
}

func (r *RedisCacheService) Delete(key ...string) {
	r.redisdb.Del(context.Background(), key...)
}

// TODO: 需要做一个测试，并于返回值
func (r *RedisCacheService) Exists(key string) bool {
	incCmd := r.redisdb.Exists(context.Background(), key)
	if incCmd.Err() != nil {
		return false
	}
	if incCmd.Val() == 0 {
		return false
	}
	return true
}

func (r *RedisCacheService) Increment(key string, val int) (int, error) {
	intCmd := r.redisdb.IncrBy(context.Background(), key, int64(val))
	if intCmd.Err() != nil {
		return int(intCmd.Val()), intCmd.Err()
	}
	return int(intCmd.Val()), nil
}

func (r *RedisCacheService) GetCache() ICache {
	return r
}
