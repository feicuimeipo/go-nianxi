package jaeger

import (
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	"go.uber.org/zap"
)

func NewConfiguration(v *viper.Viper, logger *zap.Logger) *config.Configuration {
	var (
		err error
		c   = new(config.Configuration)
	)

	if err = v.UnmarshalKey("jaeger", c); err != nil {
		logger.Sugar().Error("初始化 jaeger 配置失败:%s \n", err)
		return nil
	}
	logger.Info("加载 jaeger 配置成功")

	return c
}

func New(c *config.Configuration, logger *zap.Logger) opentracing.Tracer {

	metricsFactory := prometheus.New()
	tracer, _, err := c.NewTracer(config.Metrics(metricsFactory))

	if err != nil {
		logger.Sugar().Error("create jaeger tracer error,%v", err)
		return nil
	}

	return tracer
}

func NewJaeger(v *viper.Viper, logger *zap.Logger) opentracing.Tracer {
	c := NewConfiguration(v, logger)
	return New(c, logger)
}
