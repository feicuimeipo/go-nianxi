package test

import (
	"gitee.com/go-nianxi/go-aj-captcha/pkg/captcha/cache"
	"testing"
)

func TestLocalCacheService_Increment(t *testing.T) {
	cache := cache.NewMemCacheService(10)
	key := "test"
	cache.Increment(key, 1)

	if cache.Get(key) != "1" {
		t.Fatal("自增值不正确")
	}

	cache.Increment(key, 2)

	if cache.Get(key) != "3" {
		t.Fatal("自增值不正确")
	}
}
