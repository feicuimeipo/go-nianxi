package config

import (
	"gitee.com/go-nianxi/go-auth/internal/pkg/email"
	"gitee.com/go-nianxi/go-auth/internal/pkg/sms"
	auth_client "gitee.com/go-nianxi/go-auth/pkg/auth-client"
	client_config "gitee.com/go-nianxi/go-auth/pkg/auth-client/client-config"
	"gitee.com/go-nianxi/go-common/pkg/base"
	"gitee.com/go-nianxi/go-common/pkg/cache"
	"gitee.com/go-nianxi/go-common/pkg/captcha"
	"gitee.com/go-nianxi/go-common/pkg/db"
	"gitee.com/go-nianxi/go-common/pkg/etcd"
	"gitee.com/go-nianxi/go-common/pkg/jaeger"
	"gitee.com/go-nianxi/go-common/pkg/transports/http"
	"gitee.com/go-nianxi/go-common/pkg/valid"
	"github.com/google/wire"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
)

type AppConf struct {
	O          *AppOptions
	Logger     *zap.Logger
	Cache      cache.ICache
	Valid      *valid.Validator
	Sms        *sms.SMSClient
	Email      *email.MailClient
	Jwt        *client_config.Jwt
	AuthClient *auth_client.ClientContext
	Captcha    *captcha.AJCaptcha
}

type AppOptions struct {
	Base *base.AppOptions
	Http *http.Options `mapstructure:"-"`
	Auth *AuthOptions
}

type AuthOptions struct {
	PerDayUpdateMaxTime int                      `mapstructure:"per-day-update-max-time"`
	Domain              string                   `mapstructure:"domain"`
	Jwt                 *client_config.JwtOption `mapstructure:"jwt"`
	Dev                 *DevOptions              `mapstructure:"dev-mode"`
}

func NewSMS(v *viper.Viper, logger *zap.Logger) *sms.SMSClient {
	option, err := sms.NewOptions(v, logger)
	if err != nil {
		logger.Sugar().Error("配置信息初始化失败 %v", err)
	}
	i, err2 := sms.New(option, logger)
	if err2 != nil {
		logger.Sugar().Error("实例初始化失败 %v", err2)
		return nil
	}
	return i
}

func NewEmail(v *viper.Viper, logger *zap.Logger) *email.MailClient {
	option, err := email.NewOptions(v, logger)
	if err != nil {
		logger.Sugar().Error("配置信息初始化失败 %v", err)
	}
	email, err2 := email.New(option, logger)
	if err2 != nil {
		logger.Sugar().Error("实例初始化失败 %v", err2)
		return nil
	}
	return email
}

func NewConf(base *base.BaseConf,
	cache cache.ICache,
	valid *valid.Validator,
	sms *sms.SMSClient,
	email *email.MailClient,
	authClient *auth_client.ClientContext,
	captcha *captcha.AJCaptcha,
) *AppConf {
	Conf := new(AppConf)

	viper := base.Viper
	Conf.Logger = base.Logger

	Conf.Logger.Info("logier")
	o := new(AppOptions)
	o.Base = base.O
	o.Auth = new(AuthOptions)
	// 将读取的配置信息保存至全局变量Conf
	if err := viper.UnmarshalKey("auth", o.Auth); err != nil {
		Conf.Logger.Panic("初始化 auth 配置失败:", zap.Error(err))
		return nil
	}

	Conf.Jwt = client_config.NewJwt(o.Auth.Jwt.Secret, o.Auth.Jwt.ExpireTimeInHours, o.Auth.Jwt.ExpireTimeInHours)
	Conf.O = o
	Conf.O.Http = http.NewOption(base.Viper, base.Logger)
	Conf.Logger.Info("加载 dev-mode 配置成功")

	Conf.Cache = cache
	Conf.Valid = valid
	Conf.Sms = sms
	Conf.AuthClient = authClient
	Conf.Email = email
	Conf.Captcha = captcha
	return Conf
}

func NewLogger(baseApp *base.BaseConf) *zap.Logger {
	return baseApp.Logger
}

func NewSugaredLogger(conf *AppConf) *zap.SugaredLogger {
	return conf.Logger.Sugar()
}

type DevOptions struct {
	IsVerifyLocal bool `mapstructure:"verify-local" 		     json:"isVerifyLocal"`
}

func NewDB(v *viper.Viper, logger *zap.Logger) *gorm.DB {
	option, err := db.NewOptions(v, logger)
	if err != nil {
		logger.Sugar().Panicf("数据库配置初始化失败 %v", err)
	}
	db := db.New(option, logger)

	dbAutoMigrate(db)

	return db
}

func NewCache(v *viper.Viper, logger *zap.Logger) cache.ICache {
	option, err := cache.NewOptions(v, logger)
	if err != nil {
		logger.Sugar().Panicf("数据库配置初始化失败 %v", err)
	}
	cache, err2 := cache.New(option, logger)
	if err2 != nil {
		logger.Sugar().Panicf("数据库初始化失败 %v", err2)
	}

	return cache
}

func NewValidator(v *viper.Viper, logger *zap.Logger) *valid.Validator {
	option, err := valid.NewOptions(v, logger)
	if err != nil {
		logger.Sugar().Error("配置信息初始化失败 %v", err)
	}
	valid, err2 := valid.New(option, logger)
	if err2 != nil {
		logger.Sugar().Error("实例初始化失败 %v", err2)
		return nil
	}

	return valid
}

func NewCaptcha(v *viper.Viper, logger *zap.Logger, cache cache.ICache) *captcha.AJCaptcha {
	resource := ""
	if cache.GetRedisClient() != nil {
		return captcha.NewAJCaptchaByRDB(cache.GetRedisClient(), v, logger, resource)
	} else {
		return captcha.NewAJCaptcha(v, logger, resource)
	}
}

func NewAuthClient(v *viper.Viper) *auth_client.ClientContext {
	authClient, err := auth_client.InitNxAuthClient(base.Conf.Viper)
	if err != nil {
		log.Fatal("InitNxAuthClient+", err)
	}

	return authClient
}

func NewEtcd(v *viper.Viper, logger *zap.Logger) *clientv3.Client {
	return etcd.NewEtcd(v, logger)
}

func NewJaeger(v *viper.Viper, logger *zap.Logger) opentracing.Tracer {
	return jaeger.NewJaeger(v, logger)
}

var ProviderSet = wire.NewSet(
	base.InitViper, base.NewAppBase, NewLogger, NewSugaredLogger, NewCache, NewValidator, NewSMS, NewEmail, NewCaptcha, NewAuthClient,
	NewConf, NewDB, InitData, NewEtcd, NewJaeger)
