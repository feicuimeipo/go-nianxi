package http

import (
	"embed"
	"gitee.com/go-nianxi/go-common/pkg/transports/http/middlewares"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"html/template"
	"net/http"
	"strings"
)

const (
	StaticFSPath   = "dist"
	TemplateFSPath = "template"

	StaticUrlPrefix   = "/web"
	TemplateUrlPrefix = "/ui"

	//Index = "index"

	metricsPath = "/metrics"
	faviconPath = "/favicon.ico"
)

type Options struct {
	ContextPath    string                       `mapstructure:"context-path"      json:"contextPath"`
	Host           string                       `mapstructure:"host"              json:"host"`
	Port           int                          `mapstructure:"port"              json:"port"`
	Ssl            *SSLOptions                  `mapstructure:"ssl"               json:"ssl"`
	StaticResource *StaticResource              `mapstructure:"static-resource"   json:"port"`
	Cors           *middlewares.CorsOption      `mapstructure:"cors"              json:"cors"`
	RateLimit      *middlewares.RateLimitOption `mapstructure:"rate-limit"        json:"rateLimit"`
	EtcdDiscovery  bool                         `mapstructure:"etcd-discovery"    json:"etcdDiscovery"` //服务发布
	Mode           string                       `mapstructure:"-"                 json:"mode"`
	IsProd         bool
}

type StaticResource struct {
	Enabled        bool   `mapstructure:"enabled"                json:"enabled"`
	StaticPath     string `mapstructure:"staticPath"             json:"staticPath"`
	TemplatePath   string `mapstructure:"template-path"          json:"templatePath"`
	TemplateSuffix string `mapstructure:"template-suffix"        json:"templateSuffix"`
}

type SSLOptions struct {
	Enabled        bool   `mapstructure:"enabled"                json:"enabled"`
	Port           int    `mapstructure:"port"              json:"port"`
	CertFile       string `mapstructure:"cert-file"`
	KeyFile        string `mapstructure:"key-file"`
	TrustedCaFile  string `mapstructure:"trusted-ca-file"`
	KeyFileBytes   []byte `mapstructure:"-"                 json:"-"`
	CertFileBytes  []byte `mapstructure:"-"                 json:"-"`
	TrustedCaBytes []byte `mapstructure:"-"                 json:"-"`
}

type ext struct {
	Mode string
}

func NewOption(viper *viper.Viper, logger *zap.Logger) *Options {
	var o = new(Options)

	if err := viper.UnmarshalKey("http", o); err != nil {
		logger.Sugar().Panicf("初始化 http 配置失败:%s \n", err)
		return nil
	}

	var e = new(ext)
	if err := viper.UnmarshalKey("app", o); err != nil {
		logger.Sugar().Panicf("初始化 app.mode 配置失败:%s \n", err)
		return nil
	}

	if o.StaticResource.TemplateSuffix == "" {
		o.StaticResource.TemplateSuffix = ".html"
	}

	if o.StaticResource.TemplatePath == "" {
		o.StaticResource.TemplatePath = TemplateFSPath
	}

	o.StaticResource.TemplateSuffix = strings.TrimPrefix(o.StaticResource.TemplateSuffix, ".")

	o.Mode = e.Mode
	o.IsProd = o.Mode == "release"

	if o.Cors == nil {
		o.Cors = new(middlewares.CorsOption)
		o.Cors.Enabled = false
	}

	if o.Ssl.Enabled {
		if o.Ssl.Port == 0 {
			o.Ssl.Port = 443
		}
	}
	formatCertFile(o.Ssl)

	logger.Info("加载 http 配置成功")

	return o
}

func unescapeHTML(s string) template.HTML {
	return template.HTML(s)
}

func GinI18nLocalize(fs embed.FS) gin.HandlerFunc {
	return ginI18n.Localize(
		ginI18n.WithBundle(&ginI18n.BundleCfg{
			RootPath:         "./lang",
			AcceptLanguage:   []language.Tag{language.Chinese, language.English},
			DefaultLanguage:  language.Chinese,
			FormatBundleFile: "toml",
			UnmarshalFunc:    toml.Unmarshal,
			Loader:           &ginI18n.EmbedLoader{fs},
		}),
	)
}

func IndexRoute(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{"title": "Home Page"})
}
