package core

import (
	"gitee.com/go-nianxi/go-common/pkg/captcha/cache"
	"github.com/redis/go-redis/v9"
	"log"
	"sync"
)

// CaptchaFactory 验证码服务工厂
type CaptchaFactory struct {
	CaptchaConfig *CaptchaConfig
	CaptchaMap    map[string]CaptchaInterface
	CacheMap      map[string]cache.CacheInterface
	ServiceLock   sync.RWMutex
	CacheLock     sync.RWMutex
}

func NewCaptchaFactoryWithLocalMemory(o *CaptchaConfig) *CaptchaFactory {
	factory := &CaptchaFactory{
		CaptchaMap: make(map[string]CaptchaInterface),
		CacheMap:   make(map[string]cache.CacheInterface),
	}

	// 行为校验初始化
	var config *CaptchaConfig

	config = InitCaptchaConfig(
		cache.MemCacheKey,
		o.ResourcePath,
		o.Watermark,
		o.ClickWord,
		o.BlockPuzzle,
		2*60)

	factory.CaptchaConfig = config

	factory.RegisterCache(cache.MemCacheKey, cache.NewMemCacheService(20))

	factory.RegisterService(ClickWordCaptcha, NewClickWordCaptchaService(config.ResourcePath, config, factory.GetCache()))

	factory.RegisterService(BlockPuzzleCaptcha, NewBlockPuzzleCaptchaService(config.ResourcePath, config, factory.GetCache()))

	return factory
}

func NewCaptchaFactoryByRDB(redis redis.UniversalClient, o *CaptchaConfig) *CaptchaFactory {
	factory := &CaptchaFactory{
		CaptchaMap: make(map[string]CaptchaInterface),
		CacheMap:   make(map[string]cache.CacheInterface),
	}

	// 行为校验初始化
	var config *CaptchaConfig

	config = InitCaptchaConfig(
		cache.RedisCacheKey,
		o.ResourcePath,
		o.Watermark,
		o.ClickWord,
		o.BlockPuzzle,
		2*60)
	factory.CaptchaConfig = config
	factory.RegisterCache(cache.RedisCacheKey, cache.NewRedisCacheServiceByRDB(redis))

	factory.RegisterService(ClickWordCaptcha, NewClickWordCaptchaService(config.ResourcePath, config, factory.GetCache()))

	factory.RegisterService(BlockPuzzleCaptcha, NewBlockPuzzleCaptchaService(config.ResourcePath, config, factory.GetCache()))

	return factory
}

func (c *CaptchaFactory) GetCache() cache.CacheInterface {
	key := c.CaptchaConfig.CacheType
	c.CacheLock.RLock()
	defer c.CacheLock.RUnlock()
	if _, ok := c.CacheMap[key]; !ok {
		log.Printf("未注册%s类型的Cache", key)
	}
	return c.CacheMap[key]
}

func (c *CaptchaFactory) RegisterCache(key string, cacheInterface cache.CacheInterface) {
	c.CacheLock.Lock()
	defer c.CacheLock.Unlock()
	c.CacheMap[key] = cacheInterface
}

func (c *CaptchaFactory) RegisterService(key string, service CaptchaInterface) {
	c.ServiceLock.Lock()
	defer c.ServiceLock.Unlock()
	c.CaptchaMap[key] = service
}

func (c *CaptchaFactory) GetCaptchaService(key string) CaptchaInterface {
	c.ServiceLock.RLock()
	defer c.ServiceLock.RUnlock()
	if _, ok := c.CaptchaMap[key]; !ok {
		log.Printf("未注册%s类型的Service", key)
	}
	return c.CaptchaMap[key]
}
