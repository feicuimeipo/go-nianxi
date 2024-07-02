package grpc

import (
	"context"
	"fmt"
	"gitee.com/go-nianxi/go-common/pkg/etcd"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"time"
)

type ClientOptionalFunc func(o *ClientOptions)

type ClientOptions struct {
	Service       string        `mapstructure:"service"`
	Host          string        `mapstructure:"host"`
	Port          int           `mapstructure:"port"`
	Wait          time.Duration `mapstructure:"wait"`
	Tag           string        `mapstructure:"tag"`
	IsBlock       bool          `mapstructure:"is-block"`
	CertFile      string        `mapstructure:"cert-file"`
	KeyFile       string        `mapstructure:"key-file"`
	TrustedCaFile string        `mapstructure:"trusted-ca-file"`
	extraMeta     *ExtraMeta    `mapstructure:"extra-meta"`
	DialOptions   []grpc.DialOption
}

/*
configKey = grep.client.auth
configKey = auth.grpc.host
*/
func NewClientOptions(configKey string, v *viper.Viper, tracer opentracing.Tracer) (*ClientOptions, error) {
	var (
		err error
		o   = new(ClientOptions)
	)
	if err = v.UnmarshalKey(configKey, o); err != nil {
		return nil, fmt.Errorf("初始化 grcp.client 配置失败:%s \n", err)
	}
	if o.Service == "" {
		o.Service = fmt.Sprintf("%s:%d", o.Host, o.Port)
	}
	err = configDialOptions(o, tracer)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func configDialOptions(o *ClientOptions, tracer opentracing.Tracer) error {
	//hostname := fmt.Sprintf("%s:%d", o.Host, o.Port)
	credentials, err := tlsClient(o.CertFile, o.KeyFile, o.TrustedCaFile, o.Host)
	if err != nil {
		return err
	}
	grpc_prometheus.EnableClientHandlingTimeHistogram()
	o.DialOptions = append(o.DialOptions,
		grpc.WithTransportCredentials(credentials),
	)
	if o.extraMeta != nil {
		metaData := newCustomCredential(o.extraMeta)
		o.DialOptions = append(o.DialOptions, grpc.WithPerRPCCredentials(metaData))
	}

	if tracer != nil {
		o.DialOptions = append(o.DialOptions,
			grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
				grpc_prometheus.UnaryClientInterceptor,
				otgrpc.OpenTracingClientInterceptor(tracer)),
			),
			grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
				grpc_prometheus.StreamClientInterceptor,
				otgrpc.OpenTracingStreamClientInterceptor(tracer)),
			),
		)
	}
	return nil
}

type Client struct {
	etcdOption *etcd.Options
	o          *ClientOptions
}

func NewClient(o *ClientOptions, etcdOptions *etcd.Options, tracer opentracing.Tracer) (*Client, error) {
	configDialOptions(o, tracer)
	err := configDialOptions(o, tracer)
	if err != nil {
		return nil, err
	}

	return &Client{
		o:          o,
		etcdOption: etcdOptions,
	}, nil
}

type ClientOptional func(o *ClientOptions)

func WithGrpcDialOptions(options ...grpc.DialOption) ClientOptional {
	return func(o *ClientOptions) {
		o.DialOptions = append(o.DialOptions, options...)
	}
}

func WithTag(tag string) ClientOptional {
	return func(o *ClientOptions) {
		o.Tag = tag
	}
}

func WithTimeout(d time.Duration) ClientOptional {
	return func(o *ClientOptions) {
		o.Wait = d
	}
}

func (c *Client) Dial(options ...ClientOptionalFunc) (*grpc.ClientConn, error) {
	o := &ClientOptions{
		Wait:        c.o.Wait,
		Tag:         c.o.Tag,
		DialOptions: c.o.DialOptions,
	}

	for _, option := range options {
		option(o)
	}

	if c.etcdOption != nil && c.etcdOption.Endpoints != "" {
		if c.o.IsBlock {
			o.DialOptions = append(o.DialOptions,
				grpc.WithBlock(),
			)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		service := c.o.Service
		if service == "" {
			return nil, errors.New("service名不能为空")
		}
		target := fmt.Sprintf("consul://%s/%s?wait=%s&tag=%s", c.etcdOption.Endpoints, service, o.Wait, o.Tag)
		conn, err := grpc.DialContext(ctx, target, o.DialOptions...)
		if err != nil {
			return nil, errors.Wrap(err, "grpc dial error")
		}
		return conn, nil
	} else {
		target := fmt.Sprintf("%s:%d", c.o.Host, c.o.Port)
		conn, err := grpc.Dial(target, o.DialOptions...)
		if err != nil {
			return nil, errors.Wrap(err, "grpc dial error")
		}
		return conn, nil
	}

}
