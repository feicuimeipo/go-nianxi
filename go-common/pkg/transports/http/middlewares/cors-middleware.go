package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CorsOption struct {
	Enabled bool   `mapstructure:"enabled"       json:"enabled"`
	Origin  string `mapstructure:"origin"        json:"origin"`
}

// CORS跨域中间件
func CORSMiddleware(option *CorsOption) gin.HandlerFunc {
	//corsOption = option
	//cors.New(cors.Config{
	//	AllowOriginFunc:  func(origin string) bool { return true },
	//	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
	//	AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
	//	AllowCredentials: true,
	//	MaxAge:           12 * time.Hour,
	//})

	return func(c *gin.Context) {
		//如果需要用到跨域携带cookie，session校验的话，“Access-Control-Allow-Origin”就不能设置为“*”，要设置为具体的源地址，否则会报安全错误
		origin := c.GetHeader("Origin") //请求头部
		if origin == "" {
			if option.Origin != "" {
				c.Writer.Header().Set("Access-Control-Allow-Origin", option.Origin)
			} else {
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			}
		} else {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		//服务器支持的所有跨域请求的方法
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")

		//允许跨域设置可以返回其他子段，可以自定义字段 )
		c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length , Origin , X-Requested-With, Accept-Encoding , Access-Control-Allow-Headers, X-CSRF-Token, Authorization")

		// 允许浏览器（客户端）可以解析的头部 （重要）
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
		//设置缓存时间
		c.Header("Access-Control-Max-Age", "172800")
		//允许客户端传递校验信息比如 cookie (重要)
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			//允许类型校验
			//c.Writer.WriteHeader(http.StatusOK)
			c.AbortWithStatus(http.StatusOK)
			c.JSON(http.StatusOK, "Options Request!")
		}

		c.Next()
	}
}
