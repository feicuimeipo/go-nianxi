package client_config

import (
	"gitee.com/go-nianxi/go-auth/pkg/auth-client/auth-common"
	"gitee.com/go-nianxi/go-common/pkg/base"
	"gitee.com/go-nianxi/go-common/pkg/cache"
	nianxiGrpc "gitee.com/go-nianxi/go-common/pkg/transports/grpc"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

var Conf = new(ClientConf)
var defaultConfigFile string = "auth.xml"

type ClientConf struct {
	O      *ClientOption
	Cache  cache.ICache
	Jwt    *Jwt
	Logger *zap.Logger
}

func InitClientConf(configFile string) (*ClientConf, error) {
	if configFile == "" {
		configFile = defaultConfigFile
	}
	viper := base.InitViper(configFile)
	return InitClientAppWithViper(viper)
}

func InitClientAppWithViper(viper *viper.Viper) (*ClientConf, error) {
	client := new(ClientConf)
	o, err := newClientOptionsWithViper(viper)
	if err != nil {
		return nil, err
	}

	logger := base.NewLogger(o.Client.Log)
	client.Jwt = NewJwt(o.Jwt.Secret, o.Jwt.ExpireTimeInHours, o.Jwt.ExpireTimeInHours)
	if o.Client.Redis != nil && o.Client.Redis.DBAddress != nil {
		client.Cache = cache.NewRedis(o.Client.Redis, logger)
	} else {
		client.Cache = cache.NewLocalCacheService()
	}
	client.O = o
	client.Logger = logger
	Conf = client
	return client, nil
}

func newClientOptions(configFile string) (*ClientOption, error) {

	if configFile == "" {
		configFile = defaultConfigFile
	}

	viper := base.InitViper(configFile)
	return newClientOptionsWithViper(viper)
}

func newClientOptionsWithViper(viper *viper.Viper) (*ClientOption, error) {

	// 将读取的配置信息保存至全局变量Conf
	O := new(ClientOption)
	if err := viper.UnmarshalKey("auth", O); err != nil {
		return nil, err
	}

	if O.Client.Log == nil {
		O.Client.Log = new(base.LogOptions)
		if err := viper.UnmarshalKey("log", O.Client.Log); err != nil {
			return nil, err
		}
	}
	if O.Jwt.ExpireTimeInHours == 0 {
		O.Jwt.ExpireTimeInHours = auth_common.TOKEN_MAX_EXPIRE_HOUR
	}
	if O.Jwt.ExpireTimeInHours == 0 {
		O.Jwt.ExpireTimeInHours = auth_common.TOKEN_MAX_REFRESH_MINUTE
	}
	return O, nil
}

type ClientOption struct {
	UserType string           `mapstructure:"user-type"`
	Domain   string           `mapstructure:"domain" json:"domain"`
	Client   *lawClientOption `mapstructure:"client"`
	Jwt      *JwtOption       `mapstructure:"jwt"` //小时为单位
}

type lawClientOption struct {
	SuccessUrl string                    `mapstructure:"success-url"`
	LogoutUrl  string                    `mapstructure:"logout-url"`
	LoginUrl   string                    `mapstructure:"login-url"`
	Ignore     string                    `mapstructure:"ignore"`
	Auth       string                    `mapstructure:"auth"`
	Grpc       *nianxiGrpc.ClientOptions `mapstructure:"grpc"`
	Redis      *cache.RedisConfig        `mapstructure:"redis"`
	Log        *base.LogOptions          `mapstructure:"log"`
}

/*
 *以小时为单位
 */
func GetExpireTimeInHours() time.Duration {
	expire := Conf.O.Jwt.ExpireTimeInHours
	return time.Duration(expire) * time.Hour
}
