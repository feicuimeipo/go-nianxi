package sms

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type SMSClient struct {
}

type Options struct {
}

func (m SMSClient) SendSMS(mobile string, content string) {
	//TODO implement me
	panic("implement me")
}

func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	return &Options{}, nil
}

func New(o *Options, logger *zap.Logger) (*SMSClient, error) {
	return &SMSClient{}, nil
}

func Init(v *viper.Viper, logger *zap.Logger) *SMSClient {
	option, err := NewOptions(v, logger)
	if err != nil {
		logger.Sugar().Error("配置信息初始化失败 %v", err)
	}
	i, err2 := New(option, logger)
	if err2 != nil {
		logger.Sugar().Error("实例初始化失败 %v", err2)
		return nil
	}
	return i
}
