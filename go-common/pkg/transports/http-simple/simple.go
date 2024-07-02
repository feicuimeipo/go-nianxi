package http_simple

import (
	"embed"
	app_http "gitee.com/go-nianxi/go-common/pkg/transports/http"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func NewOption(viper *viper.Viper, logger *zap.Logger) *app_http.Options {
	return app_http.NewOption(viper, logger)
}

func NewHandlerFunc(o *app_http.Options, logger *zap.Logger, init app_http.InitRouters, initTemplate app_http.InitTemplateRouter, staticFs *embed.FS) http.Handler {
	engins := app_http.NewRouter(o, init, nil, logger)
	staticEngins := app_http.NewStaticRouter(o, initTemplate, staticFs, logger, nil, engins)
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			//如果 URL 以 /static 或 /web 开头，则走前端路由
			if strings.HasPrefix(request.URL.Path, app_http.StaticUrlPrefix) ||
				strings.HasPrefix(request.URL.Path, o.StaticResource.TemplateSuffix) ||
				strings.HasPrefix(request.URL.Path, "/") ||
				strings.HasPrefix(request.URL.Path, "/index") { // 🤔：/web 是怎么来的？
				staticEngins.Routers.ServeHTTP(writer, request)
				return
			} else {
				// 否则，走后端路由
				engins.ServeHTTP(writer, request)
				return
			}
		})
}
