package email

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type MailClient struct {
}

type Options struct {
}

func (m MailClient) SendMail(email string, title string, content string) {
	//TODO implement me
	panic("implement me")
}

func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	return &Options{}, nil
}

func New(o *Options, logger *zap.Logger) (*MailClient, error) {
	return &MailClient{}, nil
}
