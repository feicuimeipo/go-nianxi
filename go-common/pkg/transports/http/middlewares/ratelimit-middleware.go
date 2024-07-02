package middlewares

import (
	"gitee.com/go-nianxi/go-common/pkg/ecode/basic"
	"gitee.com/go-nianxi/go-common/pkg/ecode/resp"
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"net/http"
	"time"
)

var rateLimit = new(RateLimitOption)

type RateLimitOption struct {
	Enabled      bool
	FillInterval time.Duration `mapstructure:"fill-interval" json:"fillInterval"`
	Capacity     int64         `mapstructure:"capacity"      json:"capacity"`
}

// 启用限流中间件
// 默认每50毫秒填充一个令牌，最多填充200个
func RateLimitMiddleware(option *RateLimitOption) gin.HandlerFunc {
	rateLimit = option
	bucket := ratelimit.NewBucket(rateLimit.FillInterval, rateLimit.Capacity)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			resp.Fail(c, basic.RECODE_LIMITVISIT)
			c.Abort()
			return
		}
		c.Next()
	}
}

func RateLimitMiddleware1(fillInterval time.Duration, cap, quantum int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(fillInterval, cap, quantum)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			c.String(http.StatusForbidden, "rate limit...")
			c.Abort()
			return
		}
		c.Next()
	}
}
