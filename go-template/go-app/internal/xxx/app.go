package xxx

import (
	"gitee.com/go-nianxi/go-common/pkg/app"
	"gitee.com/go-nianxi/go-common/pkg/transports/grpc"
	"gitee.com/go-nianxi/go-common/pkg/transports/http"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// @Schemes:		http, https
// @BasePath		/
// @Host:			127.0.0.1
// @title			单点登录
// @version		1.0
// @description	系统项目框
// @contact.name	念小玲
// @contact.email	xlnian@163.com
func NewApp(appName string, hs *http.Server, gs *grpc.Server, logger *zap.Logger) (*app.Application, error) {

	a, err := app.NewApp(appName, logger, app.GrpcServerOption(gs), app.HttpServerOption(hs))

	if err != nil {
		return nil, errors.Wrap(err, "new app-base error")
	}

	return a, nil
}

var ProviderSet = wire.NewSet(NewApp)
