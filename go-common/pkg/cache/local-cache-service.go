package cache

import (
	"encoding/json"
	"fmt"
	"gitee.com/go-nianxi/go-common/pkg/utils"
	"github.com/coocood/freecache"
	"github.com/redis/go-redis/v9"
	"time"
)

type LocalCacheService struct {
	Cache *freecache.Cache
}

func (m *LocalCacheService) GetRedisClient() redis.UniversalClient {
	return nil
}

func (m *LocalCacheService) GetLocalCache() *freecache.Cache {
	return m.Cache
}

func NewLocalCacheService() ICache {
	cache := freecache.NewCache(100 * 1024 * 1024)
	return &LocalCacheService{Cache: cache}
}

func (m *LocalCacheService) Get(key string) (interface{}, bool) {

	val, err := m.Cache.Get([]byte(key))
	if err != nil {
		return nil, false
	}

	var ret interface{}
	err = json.Unmarshal(val, &ret)
	if err != nil {
		fmt.Sprintf("%v \n", err)
		return nil, false
	}

	return ret, true
}

func (m *LocalCacheService) GetDefault(key string, defaultValue interface{}) interface{} {
	val, found := m.Get(key)
	if !found {
		return defaultValue
	}
	return val
}

func (m *LocalCacheService) Set(key string, val interface{}, expires time.Duration) error {

	json, err := json.Marshal(val)
	if err != nil {
		return err
	}

	m.Cache.Set([]byte(key), json, int(expires.Seconds()))
	return nil
}

func (m *LocalCacheService) Delete(key ...string) {
	for _, v := range key {
		m.Cache.Del([]byte(v))
	}
}

func (m *LocalCacheService) Expire(key string, expire time.Duration) bool {
	err := m.Cache.Touch([]byte(key), int(expire.Seconds()))
	if err != nil {
		return false
	}
	return true
}

func (m *LocalCacheService) CacheType() string {
	return CACHE_MODE_LOCAL
}

func (m *LocalCacheService) Exists(key string) bool {
	_, err := m.Cache.Get([]byte(key))
	if err != nil {
		return false
	}
	return true
}

func (m *LocalCacheService) Increment(key string, val int) (int, error) {
	ret, err := m.Cache.GetOrSet([]byte(key), utils.Int2Byte(val), -1)

	if err != nil {
		return 0, fmt.Errorf("The value for %s is not an int", err)
	}

	ret2, _, err1 := m.Cache.SetAndGet([]byte(key), utils.Int2Byte(utils.Byte2Int(ret)+val), -1)
	if err1 != nil {
		return 0, fmt.Errorf("The value for %s is not an int", err1)
	}

	return utils.Byte2Int(ret2), nil
}
