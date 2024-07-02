package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
	"net/http"
)

/**
 * 基本
 */
func InitBaseRoutes(r *gin.RouterGroup) gin.IRoutes {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	//文件
	r.GET("/api-doc/*any", func(c *gin.Context) {
		swagger.DisablingWrapHandler(swaggerFiles.Handler, "SwaggerDisabled")(c)
	})

	return r

}
