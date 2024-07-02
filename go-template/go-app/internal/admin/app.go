package admin

import (
	"gitee.com/go-nianxi/go-common/pkg/app"
	"gitee.com/go-nianxi/go-common/pkg/transports/grpc"
	"gitee.com/go-nianxi/go-common/pkg/transports/http"
	"gitee.com/go-nianxi/go-template/internal/admin/repository"
	"gitee.com/go-nianxi/go-template/internal/admin/router/middleware"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// @Schemes:					http, https
// @BasePath					/
// @Host:						127.0.0.1
// @title						后台管理
// @version					1.0
// @description				系统项目框
// @contact.name				念小玲
// @contact.email				xlnian@163.com
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
func NewApplication(appName string, hs *http.Server, gs *grpc.Server, logger *zap.Logger) (*app.Application, error) {

	a, err := app.NewApp(appName, logger, app.GrpcServerOption(gs), app.HttpServerOption(hs))

	if err != nil {
		return nil, errors.Wrap(err, "new app-base error")
	}

	return a, nil
}

func WriterOperationLogChannel(dao *repository.Repository, logger *zap.Logger) error {
	// 操作日志中间件处理日志时没有将日志发送到rabbitmq或者kafka中, 而是发送到了channel中， 这里开启3个goroutine处理channel将日志记录到数据库
	logger.Info("--日志监控--")
	for i := 0; i < 3; i++ {
		go dao.SystemRepository.SaveOperationLogChannel(middleware.OperationLogChan)
	}
	return nil
}

var ProviderSet = wire.NewSet(NewApplication, WriterOperationLogChannel)
