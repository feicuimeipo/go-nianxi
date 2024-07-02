package config

import (
	"fmt"
	"gitee.com/go-nianxi/go-common/pkg/base"
)

const (
	RedisPrefixMobileVerifyCode = "mobile-verify-code" //手机验证码对应的KEY
	RedisPrefixEmailVerifyCode  = "email-verify-code"  //手机验证码对应的KEY
)

func GetCacheKey(prefix string, values []string) string {
	key := fmt.Sprintf("%s:%s", base.Conf.O.Name, prefix)
	for _, v := range values {
		if v != "" {
			key = fmt.Sprintf("%s_%s", key, v)
		}
	}
	return key
}
