package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

// RedisConfig redis配置选项
type RedisConfig struct {
	//redis单机或者集群访问地址
	Address []string
	//最大空闲连接数
	MaxIdle int
	//最大连接数
	MaxActive int
	//redis表示空闲连接保活时间
	IdleTimeout int
	//redis用户
	UserName string
	//redis密码
	PassWord string
	//是否使用redis集群
	EnableCluster bool
	//单机模式下使用redis的指定库，比如：0，1，2，3等等，默认为0
	DB int
}

type RedisUtil struct {
	Rdb redis.UniversalClient
}

func (r *RedisUtil) GetRedisClient() redis.UniversalClient {
	return r.Rdb
}

// InitConfigRedis 初始化自定义配置redis客户端（可单机， 可集群）
func NewRedisUtil(config *RedisConfig) *RedisUtil {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if config.EnableCluster {
		rdb := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    config.Address,
			Username: config.UserName,
			Password: config.PassWord,
			PoolSize: 100,
		})
		_, err := rdb.Ping(ctx).Result()
		if err != nil {
			panic(err.Error())
		}
		return &RedisUtil{Rdb: rdb}
	} else {
		rdb := redis.NewClient(&redis.Options{
			Addr:     config.Address[0],
			Username: config.UserName,
			Password: config.PassWord, // no password set
			DB:       config.DB,       // use select DB
			PoolSize: 100,             // 连接池大小
		})
		_, err := rdb.Ping(ctx).Result()
		if err != nil {
			panic(err.Error())
		}
		return &RedisUtil{Rdb: rdb}
	}
}

func NewRedisUtilByExistRDB(rdb redis.UniversalClient) *RedisUtil {
	return &RedisUtil{Rdb: rdb}
}

func (r *RedisUtil) Exists(key string) bool {
	timeVal := r.Rdb.Get(context.Background(), key+"_HoldTime").Val()
	cacheHoldTime, err := strconv.ParseInt(timeVal, 10, 64)

	if err != nil {
		return false
	}

	if cacheHoldTime == 0 {
		return true
	}

	if cacheHoldTime < time.Now().Unix() {
		r.Delete(key)
		return false
	}
	return true
}

func (r *RedisUtil) Get(key string) string {
	val := r.Rdb.Get(context.Background(), key).Val()
	return val
}

func (r *RedisUtil) Set(key string, val string, expiresInSeconds int) {
	//设置阈值，达到即clear缓存
	rdsResult := r.Rdb.Set(context.Background(), key, val, time.Duration(expiresInSeconds)*time.Second)
	fmt.Println("rdsResult: ", rdsResult.String(), "rdsErr: ", rdsResult.Err())
	if expiresInSeconds > 0 {
		// 缓存失效时间
		nowTime := time.Now().Unix() + int64(expiresInSeconds)
		r.Rdb.Set(context.Background(), key+"_HoldTime", strconv.FormatInt(nowTime, 10), time.Duration(expiresInSeconds)*time.Second)
	} else {
		r.Rdb.Set(context.Background(), key+"_HoldTime", strconv.FormatInt(0, 10), time.Duration(expiresInSeconds)*time.Second)
	}
}

func (r *RedisUtil) Delete(key string) {
	r.Rdb.Del(context.Background(), key)
	r.Rdb.Del(context.Background(), key+"_HoldTime")
}

func (l *RedisUtil) Clear() {
	//for key, _ := range r.Data {
	//	r.Delete(key)
	//}
}
