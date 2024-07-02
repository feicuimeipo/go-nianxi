package admin

import (
	"gitee.com/go-nianxi/go-common/pkg/app"
	"gitee.com/go-nianxi/go-common/pkg/transports/grpc"
	"gitee.com/go-nianxi/go-common/pkg/transports/http"
	"gitee.com/go-nianxi/go-admin/internal/admin/dao"
	"gitee.com/go-nianxi/go-admin/internal/admin/router/middleware"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func NewApplication(appName string, hs *http.Server, gs *grpc.Server, logger *zap.Logger) (*app.Application, error) {

	a, err := app.NewApp(appName, logger, app.GrpcServerOption(gs), app.HttpServerOption(hs))

	if err != nil {
		return nil, errors.Wrap(err, "new app-base error")
	}

	return a, nil
}

func WriterOperationLogChannel(dao *dao.Dao, logger *zap.Logger) error {
	// 操作日志中间件处理日志时没有将日志发送到rabbitmq或者kafka中, 而是发送到了channel中， 这里开启3个goroutine处理channel将日志记录到数据库
	logger.Info("--日志监控--")
	for i := 0; i < 3; i++ {
		go dao.OperationLogDao.SaveOperationLogChannel(middleware.OperationLogChan)
	}
	return nil
}

var ProviderSet = wire.NewSet(NewApplication, WriterOperationLogChannel)
