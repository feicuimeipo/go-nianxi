package middleware

import (
	"gitee.com/go-nianxi/go-auth/internal/auth/model"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

/**
 * 操作日志channel
 */
var LoginAuditLogChan = make(chan *model.LoginAuditLog, 30)

func (middleware *Middleware) OperationLogMiddleware(contextPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行耗时
		timeCost := endTime.Sub(startTime).Milliseconds()

		// 获取当前登录用户
		var username string
		ctxUser, exists := c.Get("user")
		if !exists {
			username = "未登录"
		}
		user, ok := ctxUser.(model.User)
		if !ok {
			username = "未登录"
		}
		username = user.Username

		// 获取访问路径
		path := strings.TrimPrefix(c.FullPath(), "/"+contextPath)
		// 请求方式
		method := c.Request.Method

		// 获取接口描述

		loginAuditLog := model.LoginAuditLog{
			Username:   username,
			Ip:         c.ClientIP(),
			IpLocation: "",
			Method:     method,
			Path:       path,
			Status:     c.Writer.Status(),
			StartTime:  startTime,
			TimeCost:   timeCost,
			//UserAgent:  c.Request.UserAgent(),
		}

		// 最好是将日志发送到rabbitmq或者kafka中
		// 这里是发送到channel中，开启3个goroutine处理
		LoginAuditLogChan <- &loginAuditLog
	}
}
