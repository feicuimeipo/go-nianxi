package cache

const (
	// MemCacheKey 内存缓存标识
	MemCacheKey = "mem"
	// RedisCacheKey redis缓存标识
	RedisCacheKey = "redis"
)

type CacheInterface interface {
	Get(key string) string
	Set(key string, val string, expiresInSeconds int)
	Delete(key string)
	Exists(key string) bool
	GetType() string
	Increment(key string, val int) int
}
