package captcha

import (
	"gitee.com/go-nianxi/go-common/pkg/captcha/core"
	"gitee.com/go-nianxi/go-common/pkg/captcha/router"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type AJCaptcha struct {
	O          *core.CaptchaConfig
	Factory    *core.CaptchaFactory
	HandleFunc *router.HandleFunc
}

func NewOptions(v *viper.Viper, logger *zap.Logger) (*core.CaptchaConfig, error) {
	var (
		o = new(core.CaptchaConfig)
	)
	if v == nil {
		err := v.UnmarshalKey("captcha", o)
		if err != nil {
			logger.Sugar().Error("初始化 cahe 配置失败:%s \n", err)
			return nil, err
		}

		logger.Info("加载 cache 配置成功")

	}
	return o, nil
}

func NewAJCaptcha(v *viper.Viper, logger *zap.Logger, resourceRootPath string) *AJCaptcha {
	o, err := NewOptions(v, logger)
	if err != nil {
		logger.Sugar().Error("验证码加载出错！%v", err)
		return nil
	}
	logger.Sugar().Info("o.ResourcePath,%s", o.ResourcePath)

	if resourceRootPath != "" {
		o.ResourcePath = resourceRootPath
	}
	factory := core.NewCaptchaFactoryWithLocalMemory(o)
	HandleFunc := router.NewHandleFunc(factory)
	return &AJCaptcha{
		O:          o,
		HandleFunc: HandleFunc,
		Factory:    factory,
	}
}

func NewAJCaptchaByRDB(RDB redis.UniversalClient, v *viper.Viper, logger *zap.Logger, resourceRootPath string) *AJCaptcha {
	o, err := NewOptions(v, logger)
	if err != nil {
		logger.Sugar().Error("验证码加载出错！%v", err)
		return nil
	}
	if resourceRootPath != "" {
		o.ResourcePath = resourceRootPath
	}

	logger.Sugar().Info("o.ResourcePath,%s", o.ResourcePath)

	var factory *core.CaptchaFactory
	if RDB == nil {
		factory = core.NewCaptchaFactoryWithLocalMemory(o)
	} else {
		factory = core.NewCaptchaFactoryByRDB(RDB, o)
	}

	HandleFunc := router.NewHandleFunc(factory)
	return &AJCaptcha{
		O:          o,
		HandleFunc: HandleFunc,
		Factory:    factory,
	}
}
