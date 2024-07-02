package auth_common

import "fmt"

const (
	TOKEN_MAX_REFRESH_MINUTE = 15 // token还有多久过期就返回新token
	TOKEN_MAX_EXPIRE_HOUR    = 24 //token最长有效期
)

const (
	TokenHeaderNamePrefix = "Bearer " // header名称

	NewTokenName = "new-token"

	TokenHeaderName = "Authorization"
	TokenQueryName  = "token"

	TokenCookiesName = "nx-token"

	TokenNameLookup = "header: " + TokenHeaderName + ", query: " + TokenQueryName + ", cookie: " + TokenCookiesName // 自动在这几个地方寻找请求中的token

	ContextAuthUser  = "ctx-auth-user-key"
	ContextTokenName = "ctx-token-key"
)

const (
	FLAG_LOOGIN       = "login" //登录跳转
	FLAG_SUCCESS_SALT = "nianxi_auth"
)
const (
	X_Requested_With_HEADER = "X-Requested-With"
	X_Requested_With_VALUE  = "XMLHttpRequest"
)

func GetAuthCacheKey(uid uint) string {
	return fmt.Sprintf("%s:%d", "nx-auth", uid)
}

func GetAuthCacheKeyString(id string) string {
	return fmt.Sprintf("%s:%s", "nx-auth", id)
}
