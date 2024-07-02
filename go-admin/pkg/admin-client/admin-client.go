package admin_client

import (
	api_impl "gitee.com/go-nianxi/go-admin/pkg/admin-client/apimpl"
	client_config "gitee.com/go-nianxi/go-admin/pkg/admin-client/config"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var _clientContext = new(ClientContext)

type ClientContext struct {
	apiImpl *api_impl.ApiImpl
	conf    *client_config.ClientConf
}

func NewClientContext(configFile string) (*ClientContext, error) {

	conf, err := client_config.InitClientConf(configFile)
	if err != nil {
		return nil, errors.Wrap(err, "客户端配置出错！")
	}
	_clientContext.conf = conf

	_clientContext.apiImpl = api_impl.NewApiImpl(conf.O.Grpc, conf.Logger)
	return _clientContext, nil
}

func NewClientContextWithViper(viper *viper.Viper) (*ClientContext, error) {
	var err error

	conf, err := client_config.InitClientAppWithViper(viper)
	if err != nil {
		return nil, errors.Wrap(err, "客户端配置出错！")
	}
	_clientContext.conf = conf

	_clientContext.apiImpl = api_impl.NewApiImpl(conf.O.Grpc, conf.Logger)
	return _clientContext, nil
}
