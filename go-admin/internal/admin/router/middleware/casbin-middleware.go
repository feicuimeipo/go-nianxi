package middleware

import (
	"gitee.com/go-nianxi/go-common/pkg/base"
	"gitee.com/go-nianxi/go-common/pkg/ecode/resp"
	"gitee.com/go-nianxi/go-admin/internal/admin/model"
	"gitee.com/go-nianxi/go-admin/internal/admin/util"
	"github.com/gin-gonic/gin"
	"strings"
	"sync"
)

var checkLock sync.Mutex

// Casbin中间件, 基于RBAC的权限访问控制模型
func (m *Middleware) CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exist := c.Get("user")
		if !exist {
			var err error
			var claim map[string]interface{}
			var user model.User
			claim, err = m.AuthMiddleware.GetClaimsFromJWT(c)
			if err == nil {
				if claimUser, ok := claim["user"].(model.User); ok {
					user = claimUser
					c.Set("user", user)
				} else if userStr, ok1 := claim["user"].(string); ok1 {
					util.Json2Struct(userStr, &user)
					c.Set("user", user)
				} else {
					userMap := claim["user"].(map[string]interface{})
					userStr = util.Struct2Json(userMap)
					if user.ID != 0 && user.Username != "" {
						util.Json2Struct(userStr, &user)
						c.Set("user", user)
					}
				}

			}
		}
		user, err := m.userDao.GetCurrentUser(c)
		if err != nil {
			resp.Writer(c, 401, 401, nil, err.Error())
			//response.Response(c, 401, 401, nil, "用户未登录")
			c.Abort()
			return
		}
		if user.Status != 1 {
			resp.Writer(c, 402, 401, nil, "当前用户已被禁用")
			//response.Response(c, 401, 401, nil, "当前用户已被禁用")
			c.Abort()
			return
		}
		// 获得用户的全部角色
		roles := user.Roles
		// 获得用户全部未被禁用的角色的Keyword
		var subs []string
		for _, role := range roles {
			if role.Status == 1 {
				subs = append(subs, role.Keyword)
			}
		}
		// 获得请求路径URL
		// 获取请求方式
		obj := strings.TrimPrefix(c.FullPath(), base.UrlPathPrefix)
		baseUrl := c.Request.Host
		baseUrl = strings.TrimPrefix(baseUrl, "http://")
		baseUrl = strings.TrimPrefix(baseUrl, "https://")
		act := c.Request.Method
		isPass := m.check(subs, obj, act, baseUrl)
		if !isPass {
			resp.Writer(c, 402, 402, nil, "没有权限")
			//response.Response(c, 401, 401, nil, "没有权限")
			c.Abort()
			return
		}
		c.Next()
	}
}
func (m *Middleware) check(subs []string, obj string, act string, url string) bool {
	if !strings.HasPrefix(obj, "/") {
		obj = "/" + obj
	}
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")
	// 同一时间只允许一个请求执行校验, 否则可能会校验失败
	checkLock.Lock()
	defer checkLock.Unlock()
	isPass := false
	for _, sub := range subs {
		pass, _ := m.casbin.Enforce(sub, obj, act, url)
		if pass {
			isPass = true
			break
		}
	}
	return isPass
}
