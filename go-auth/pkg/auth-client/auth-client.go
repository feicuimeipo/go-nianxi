package auth_client

import (
	"gitee.com/go-nianxi/go-auth/pkg/auth-client/apimpl"
	"gitee.com/go-nianxi/go-auth/pkg/auth-client/client-config"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type ClientContext struct {
	apiImpl *apimpl.ApiImpl
	conf    *client_config.ClientConf
}

func InitNxAuthClientByConfigFile(configFile string) (*ClientContext, error) {
	var _clientContext = new(ClientContext)

	conf, err := client_config.InitClientConf(configFile)
	if err != nil {
		return nil, errors.Wrap(err, "客户端配置出错！")
	}
	_clientContext.conf = conf

	_clientContext.apiImpl = apimpl.NewApiImpl(conf.O.Client.Grpc, conf.Logger)
	return _clientContext, nil
}

func InitNxAuthClient(viper *viper.Viper) (*ClientContext, error) {
	var err error

	var _clientContext = new(ClientContext)

	conf, err := client_config.InitClientAppWithViper(viper)
	if err != nil {
		return nil, errors.Wrap(err, "客户端配置出错！")
	}
	_clientContext.conf = conf
	_clientContext.apiImpl = apimpl.NewApiImpl(conf.O.Client.Grpc, conf.Logger)
	return _clientContext, nil
}
