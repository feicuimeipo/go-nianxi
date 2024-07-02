package middleware

import (
	"fmt"
	"gitee.com/go-nianxi/go-common/pkg/base"
	"gitee.com/go-nianxi/go-common/pkg/ecode/resp"
	"gitee.com/go-nianxi/go-common/pkg/transports/http"
	"gitee.com/go-nianxi/go-admin/internal/admin/config"
	"gitee.com/go-nianxi/go-admin/internal/admin/domain/vo"
	"gitee.com/go-nianxi/go-admin/internal/admin/model"
	"gitee.com/go-nianxi/go-admin/internal/admin/util"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

// 初始化jwt中间件
func (m *Middleware) InitAuth() (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           config.Conf.O.Jwt.Realm,                                          // jwt标识
		Key:             []byte(config.Conf.O.Jwt.Secret),                                 // 服务端密钥
		Timeout:         time.Hour * time.Duration(config.Conf.O.Jwt.ExpireTimeInHours),   // token过期时间
		MaxRefresh:      time.Hour * time.Duration(config.Conf.O.Jwt.MaxRefreshInMinutes), // token最大刷新时间(RefreshToken过期时间=Timeout+MaxRefresh)
		PayloadFunc:     m.payloadFunc,                                                    // 有效载荷处理
		IdentityHandler: m.identityHandler,                                                // 解析Claims
		Authenticator:   m.login,                                                          // 校验token的正确性, 处理登录逻辑
		Unauthorized:    m.unauthorized,                                                   // 用户登录校验失败处理
		LoginResponse:   m.loginResponse,                                                  // 登录成功后的响应
		LogoutResponse:  m.logoutResponse,                                                 // 登出后的响应
		RefreshResponse: m.refreshResponse,                                                // 刷新token后的响应
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",               // 自动在这几个地方寻找请求中的token
		TokenHeadName:   "Bearer ",                                                        // header名称
		TimeFunc:        time.Now,
	})
	return authMiddleware, err
}

// 有效载荷处理
func (m *Middleware) payloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(map[string]interface{}); ok {
		var user model.User
		// 将用户json转为结构体
		util.JsonI2Struct(v["user"], &user)
		return jwt.MapClaims{
			jwt.IdentityKey: user.ID,
			"user":          v["user"],
		}
	}
	return jwt.MapClaims{}
}

// 解析Claims
func (m *Middleware) identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	//此处返回值类型map[string]interface{}与payloadFunc和authorizator的data类型必须一致, 否则会导致授权失败还不容易找到原因
	return map[string]interface{}{
		"IdentityKey": claims[jwt.IdentityKey],
		"user":        claims["user"],
	}
}

// 用户登录校验成功处理
func (m *Middleware) authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(map[string]interface{}); ok {
		userStr := v["user"].(string)
		var user model.User
		// 将用户json转为结构体
		util.Json2Struct(userStr, &user)
		// 将用户保存到context, api调用时取数据方便
		c.Set("user", user)
		return true
	}
	return false
}

// 校验token的正确性, 处理登录逻辑
//
//	@Tags		认证
//	@summary	登录-帐号与密码
//	@Accept		json
//	@Produce	json
//	@Response	200	{object}	resp.ResponseMsg
//	@Param		req	body		vo.RegisterAndLoginRequest	true	"登录信息"
//	@Router		/auth/login [post]
func (m *Middleware) login(c *gin.Context) (interface{}, error) {
	var req vo.RegisterAndLoginRequest
	// 请求json绑定
	if err := c.ShouldBind(&req); err != nil {
		return "", err
	}

	u := &model.User{Username: req.Username}

	// 密码通过RSA解密
	if m.httOption.Mode == base.PROD_MODE || len(req.Password) > 20 {
		decodeData, err := http.RSADecrypt([]byte(req.Password), m.httOption.Ssl.KeyFileBytes)
		if err != nil {
			return nil, err
		}
		u.Password = string(decodeData)
	} else {
		u.Password = req.Password
	}

	// 密码校验
	user, err := m.userDao.Login(u)
	if err != nil {
		return nil, err
	}
	userStr := util.Struct2Json(*user)
	// 将用户以json格式写入, payloadFunc/authorizator会使用到
	return map[string]interface{}{
		"user": userStr,
	}, nil
}

// 用户登录校验失败处理
func (m *Middleware) unauthorized(c *gin.Context, code int, message string) {
	config.Conf.Logger.Error("JWT认证失败, 错误码: %d, 错误信息: %s", zap.Int("code", code), zap.String("message", message))
	resp.Writer(c, code, code, nil, fmt.Sprintf("JWT认证失败, 错误码: %d, 错误信息: %s", code, message))
}

// 登录成功后的响应
func (m *Middleware) loginResponse(c *gin.Context, code int, token string, expires time.Time) {
	//response.Response(c, code, code,
	//	gin.H{
	//		"token":   token,
	//		"expires": expires.Format("2006-01-02 15:04:05"),
	//	},
	//	"登录成功")
	resp.Writer(c, code, code, gin.H{
		"token":   token,
		"expires": expires.Format("2006-01-02 15:04:05"),
	}, "登录成功")
}

// 登出后的响应
func (m *Middleware) logoutResponse(c *gin.Context, code int) {
	resp.OK(c, nil, "退出成功")
}

// 刷新token后的响应
func (m *Middleware) refreshResponse(c *gin.Context, code int, token string, expires time.Time) {
	resp.Writer(c, code, code, gin.H{
		"token":   token,
		"expires": expires,
	},
		"刷新token成功")
}
