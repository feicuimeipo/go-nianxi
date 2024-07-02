package etcd

import (
	"crypto/tls"
	"fmt"
	"github.com/spf13/viper"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"strings"
	"time"
)

type Options struct {
	Endpoints     string        `mapstructure:"endpoints"        json:"endpoints"`      //端口 >999开启，否则关关地
	Username      string        `mapstructure:"username"   json:"username"`             //端口 为空，不开启
	Password      string        `mapstructure:"password"   json:"password"`             //端口 为空，不开启
	DialTimeout   time.Duration `mapstructure:"dial-timeout"   json:"dialTimeout"`      //端口 为空，不开启
	CertFile      string        `mapstructure:"certs-file"   json:"certFile"`           //证书
	KeyFile       string        `mapstructure:"key-file"   json:"keyFile"`              //秘钥
	TrustedCAFile string        `mapstructure:"trusted-ca-file"   json:"trustedCAFile"` //ca证书
}

func NewOptions(v *viper.Viper, logger *zap.Logger) *Options {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("etcd", o); err != nil {
		logger.Sugar().Panicf("初始化 etcd 配置信息失败:%s \n", err)
		return nil
	}

	return o
}

func NewSimple(address string) *clientv3.Client {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(address, ";"),
		DialTimeout: 15 * time.Second,
	})
	if err != nil {
		fmt.Printf("连接etcdOptions失败：%s\n", err)
		return nil
	}

	return cli
}

func NewEtcd(v *viper.Viper, logger *zap.Logger) *clientv3.Client {
	o := NewOptions(v, logger)
	return New(o, logger)
}

func New(o *Options, logger *zap.Logger) *clientv3.Client {
	//var endpoint = config.Conf.EtcdOptions.Endpoints
	if o.Endpoints == "" {
		logger.Error("endpoints为空！")
		return nil
	}
	endpoints := strings.Split(o.Endpoints, ",; ")
	var err error
	var tlsConfig *tls.Config

	if o.CertFile != "" {
		tlsInfo := transport.TLSInfo{
			CertFile:      o.CertFile,      // etcdOptions公钥
			KeyFile:       o.KeyFile,       // etcdOptions私钥
			TrustedCAFile: o.TrustedCAFile, // ca证书
		}
		tlsConfig, err = tlsInfo.ClientConfig()
		if err != nil {
			logger.Sugar().Fatal(err)
		}
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,     // etcdOptions集群列表
		DialTimeout: o.DialTimeout, // 连接超时时间
		TLS:         tlsConfig,     //不使用tls可不用添加
		Username:    o.Username,    //不开启权限认证可不用添加，password同
		Password:    o.Password,
	})
	if err != nil {
		logger.Sugar().Error(err)
		return nil
	}

	logger.Info("加载 etcd 配置成功")
	return cli
}
