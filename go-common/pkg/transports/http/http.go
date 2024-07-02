package http

import (
	"context"
	"embed"
	"fmt"
	"gitee.com/go-nianxi/go-common/pkg/transports/http/middlewares"
	"gitee.com/go-nianxi/go-common/pkg/transports/http/middlewares/ginprom"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-contrib/pprof" //性能分析工具
	"github.com/gin-contrib/zap"   //日志
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"html/template"
	"io/fs"
	"net/http"
	"strings"
	"time"
)

type StaticEngine struct {
	Routers  *gin.Engine
	StaticFs *embed.FS
}

type Server struct {
	o            *Options
	app          string
	port         int
	host         string
	ssl          *SSLOptions
	logger       *zap.Logger
	routers      *gin.Engine
	staticRouter *StaticEngine
	httpServer   *http.Server
	discovery    IDiscovery
}

type IDiscovery interface {
	register(addr string, app string) error
	deRegister(addr string, app string) error
	client() interface{}
}

type InitRouters func(router *gin.RouterGroup, engine *gin.Engine)

func NewRouter(o *Options, initRouters InitRouters, tracer opentracing.Tracer, logger *zap.Logger) *gin.Engine {
	if o != nil {
		gin.SetMode(o.Mode)
		gin.ForceConsoleColor()

	}

	// 配置gin
	engine := gin.New()
	engine.Use(gin.Recovery()) // panic之后自动恢复
	engine.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	engine.Use(ginzap.RecoveryWithZap(logger, true))
	engine.Use(ginprom.New(engine).Middleware()) // 添加prometheus 监控
	// 初始化JWT认证中间件
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	if o.Cors.Enabled {
		engine.Use(middlewares.CORSMiddleware(o.Cors))
	}
	if o.RateLimit != nil && o.RateLimit.Enabled {
		engine.Use(middlewares.RateLimitMiddleware(o.RateLimit))
	}

	pprof.Register(engine)

	// 路由分组
	apiGroup := engine.Group("/")
	apiGroup.GET(metricsPath, gin.WrapH(promhttp.Handler()))
	if initRouters != nil {
		initRouters(apiGroup, engine)
	}

	return engine
}

func NewStaticRouter(o *Options, staticFs *embed.FS, logger *zap.Logger, tracer opentracing.Tracer, engine *gin.Engine) *StaticEngine {
	if !o.StaticResource.Enabled {
		return nil
	}
	// 配置gin
	if engine == nil {
		gin.SetMode(o.Mode)
		engine = gin.New()
		engine.Use(gin.Recovery())
		pprof.Register(engine)
	}

	//静态资源
	static, err1 := fs.Sub(staticFs, StaticFSPath)
	if err1 != nil {
		logger.Sugar().Panicf("ginEngine failed to config static: %v", err1)
	}
	engine.StaticFS(StaticUrlPrefix, http.FS(static))
	engine.StaticFile("/favicon.ico", "./favicon.ico")
	engine.StaticFile("/robots.txt", "./robots.txt")
	funcMap := template.FuncMap{
		"UnescapeHTML": unescapeHTML,
		"Localize":     ginI18n.GetMessage,
	}
	engine.FuncMap = funcMap

	return &StaticEngine{
		Routers:  engine,
		StaticFs: staticFs,
	}

}

func NewServer(o *Options, routers *gin.Engine, staticRouter *StaticEngine, etcClient *clientv3.Client, logger *zap.Logger) *Server {
	var discovery IDiscovery
	if etcClient != nil {
		discovery = newEtcdServer(5, etcClient)
	} else {
		discovery = nil
	}

	s := &Server{
		logger:       logger.With(zap.String("type", "http.Server")),
		routers:      routers,
		o:            o,
		discovery:    discovery,
		staticRouter: staticRouter,
	}

	return s
}

func (s *Server) Application(name string) {
	s.app = name
}

func (s *Server) handlerFunc() http.HandlerFunc {
	handfunc := http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			if strings.HasPrefix(request.URL.Path, StaticUrlPrefix) ||
				strings.HasPrefix(request.URL.Path, TemplateUrlPrefix) { // 🤔：/web 是怎么来的？
				s.staticRouter.Routers.ServeHTTP(writer, request)
			} else {
				// 否则，走后端路由
				s.routers.ServeHTTP(writer, request)
			}
		})
	return handfunc
}

func (s *Server) Start() error {
	s.port = s.o.Port
	s.ssl = s.o.Ssl
	s.host = s.o.Host
	if s.port == 0 && s.ssl.Port == 0 {
		return errors.New("端口不可以为空！") //s.port = netutil.GetAvailablePort()
	}

	if s.discovery != nil {
		if s.host == "" {
			return errors.New("Host地址不可以为空！")
		}
	}

	if !s.ssl.Enabled {
		addr := fmt.Sprintf("%s:%d", "", s.port)
		s.httpServer = &http.Server{Addr: addr, Handler: s.handlerFunc()}

		s.logger.Info("http server starting...", zap.String("addr", addr))
		go func() {
			if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				s.logger.Fatal("服务启动失败:", zap.Error(err))
				return
			}
		}()

		if s.discovery != nil {
			addr = fmt.Sprintf("%s:%d", s.host, s.port)
			if err := s.discovery.register(addr, s.app); err != nil {
				return errors.New("register http server error")
			}
		}
	} else {
		formatCertFile(s.ssl)
		addr := fmt.Sprintf("%s:%d", "", s.ssl.Port)

		s.httpServer = &http.Server{Addr: addr, Handler: s.handlerFunc()}

		s.logger.Info("http server starting(SSL Enabled)...", zap.String("addr", addr))
		go func() {
			if err := s.httpServer.ListenAndServeTLS(s.ssl.CertFile, s.ssl.KeyFile); err != nil && err != http.ErrServerClosed {
				s.logger.Fatal("服务启动失败:", zap.Error(err))
				return
			}
		}()

		if s.discovery != nil {
			addr = fmt.Sprintf("%s:%d", s.host, s.ssl.Port)
			if err := s.discovery.register(addr, s.app+"_ssl"); err != nil {
				return errors.New("register http server error")
			}
		}
	}

	return nil
}

func (s *Server) Stop() error {
	s.logger.Info("==http服务启动停止...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5) // 平滑关闭,等待5秒钟处理
	defer cancel()

	if s.discovery != nil {
		if !(s.ssl.Enabled) {
			addr := fmt.Sprintf("%s:%d", s.host, s.port)
			if err := s.discovery.deRegister(addr, s.app); err != nil {
				return errors.Wrap(err, "deregister http server error")
			}
		} else {
			addr := fmt.Sprintf("%s:%d", s.host, s.ssl.Port)
			if err := s.discovery.deRegister(addr, s.app+"_ssl"); err != nil {
				return errors.Wrap(err, "deregister http server error")
			}
		}
	}

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "shutdown http server error")
	}

	return nil
}
