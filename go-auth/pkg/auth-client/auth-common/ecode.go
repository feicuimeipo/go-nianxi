package auth_common

const (
	SUCCESS                   = 0
	ERROR_AHTHORIZATION_EMPTY = 1402
	ERROR_TOKEN_EMPTY         = 1403
	ERROR_TOKEN_VALIDATE      = 1404
	ERROR_TOKEN_EXPIRED       = 1405
	ERROR_SIGN_VALIDATE       = 1406
	ERROR_SIGN_NOT_NULL       = 1407
)

var errMap = map[int]string{
	SUCCESS:                   "成功",
	ERROR_AHTHORIZATION_EMPTY: "授权码 为空 请重新登录",
	ERROR_TOKEN_EMPTY:         "token 为空 请重新登录",
	ERROR_TOKEN_VALIDATE:      "Token 无效，请重新登录",
	ERROR_TOKEN_EXPIRED:       "Token 已过期，请重新登录",
	ERROR_SIGN_VALIDATE:       "签名无效",
	ERROR_SIGN_NOT_NULL:       "签名不能为空",
}

func GetCodeMsg(code int) string {
	return errMap[code]
}
