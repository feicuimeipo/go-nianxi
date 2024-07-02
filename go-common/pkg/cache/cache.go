package cache

import (
	"github.com/coocood/freecache"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"time"
)

type ICache interface {
	Set(key string, val interface{}, expire time.Duration) error
	GetDefault(key string, defaultValue interface{}) interface{}
	Get(key string) (interface{}, bool)
	Delete(key ...string)
	CacheType() string
	Increment(key string, val int) (int, error)
	GetRedisClient() redis.UniversalClient
	GetLocalCache() *freecache.Cache
	Expire(key string, expire time.Duration) bool
	Exists(key string) bool
}

const (
	// For use with functions that take an expiration time.
	NoExpiration time.Duration = -1
	// For use with functions that take an expiration time. Equivalent to
	// passing in the same expiration duration as was given to New() or
	// NewFrom() when the datacache was created (e.g. 5 minutes.)
	DefaultExpiration time.Duration = 0
)

const (
	CACHE_MODE_LOCAL = "local"
	CACHE_MODE_REDIS = "redis"
	CACHE_MODE_BOTH  = "both"
)

type Options struct {
	CacheMode string       `mapstructure:"mode"`
	Redis     *RedisConfig `mapstructure:"redis"`
}

type RedisConfig struct {
	DBAddress     []string `mapstructure:"db-address" json:"dbAddress"`          //redis单机或者集群访问地址
	DBMaxIdle     int      `mapstructure:"db-max-idle" json:"dbMaxIdle"`         //最大空闲连接数
	DBMaxActive   int      `mapstructure:"db-max-active" json:"dbMaxActive"`     //最大连接数
	DBIdleTimeout int      `mapstructure:"db-idle-timeout" json:"dbIdleTimeout"` //redis表示空闲连接保活时间
	DBUserName    string   `mapstructure:"db-user-name" json:"dbUserName"`       //redis用户
	DBPassWord    string   `mapstructure:"db-password" json:"dbPassWord"`        //redis密码
	EnableCluster bool     `mapstructure:"enable-cluster" json:"enableCluster"`  //是否使用redis集群
	DB            int      `mapstructure:"db" json:"db"`                         //单机模式下使用redis的指定库，比如：0，1，2，3等等，默认为0
}

func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		cache = new(Options)
	)
	err := v.UnmarshalKey("cache", cache)
	if err != nil {
		logger.Sugar().Panicf("初始化 cahe 配置失败:%s \n", err)
		return nil, err
	}
	logger.Info("加载 cache 配置成功")

	return cache, nil
}

func New(Options *Options, logger *zap.Logger) (ICache, error) {
	var cacheMap = make(map[string]ICache)
	cacheMode := Options.CacheMode
	if cacheMode == CACHE_MODE_BOTH || cacheMode == CACHE_MODE_REDIS {
		cacheMap[CACHE_MODE_REDIS] = NewRedis(Options.Redis, logger)
		return cacheMap[CACHE_MODE_REDIS], nil
	}
	if cacheMode == CACHE_MODE_LOCAL || cacheMode == CACHE_MODE_REDIS {
		cacheMap[CACHE_MODE_LOCAL] = NewLocalCacheService()
		return cacheMap[CACHE_MODE_LOCAL], nil
	}
	return nil, nil
}

func NewLocalCache() *freecache.Cache {
	cache := freecache.NewCache(100 * 1024 * 1024)
	return cache
}

func NewRedisClient(config *RedisConfig, logger *zap.Logger) redis.UniversalClient {
	var redisdb redis.UniversalClient
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if config.EnableCluster {
		redisdb = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    config.DBAddress,
			Username: config.DBUserName,
			Password: config.DBPassWord,
			PoolSize: 100,
		})
		_, err := redisdb.Ping(ctx).Result()
		if err != nil {
			panic(err.Error())
		}
	} else {
		redisdb = redis.NewClient(&redis.Options{
			Addr:     config.DBAddress[0],
			Username: config.DBUserName,
			Password: config.DBPassWord,
			DB:       config.DB, // use select DB
			PoolSize: 100,       // 连接池大小
		})
		_, err := redisdb.Ping(ctx).Result()
		if err != nil {
			panic(err.Error())
		}
	}
	logger.Sugar().Info("Redis成功配置")
	return redisdb
}
