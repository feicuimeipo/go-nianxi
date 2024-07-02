package app

import (
	"gitee.com/go-nianxi/go-common/pkg/base"
	"gitee.com/go-nianxi/go-common/pkg/transports/grpc"
	"gitee.com/go-nianxi/go-common/pkg/transports/http"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

type Application struct {
	base.AppOptions
	Name       string
	logger     *zap.Logger
	httpServer *http.Server
	grpcServer *grpc.Server
}

type Components func(app *Application) error

func HttpServerOption(svr *http.Server) Components {
	return func(app *Application) error {
		svr.Application(app.Name)
		app.httpServer = svr
		return nil
	}
}

func GrpcServerOption(svr *grpc.Server) Components {
	return func(app *Application) error {
		svr.Application(app.Name)
		app.grpcServer = svr
		return nil
	}
}

// name=""则从配置文件中取
func NewApp(name string, logger *zap.Logger, options ...Components) (*Application, error) {
	app := &Application{
		logger: logger.With(zap.String("type", "Application")),
	}
	if name != "" {
		app.Name = name
	}

	for _, option := range options {
		if err := option(app); err != nil {
			return nil, err
		}
	}

	return app, nil
}

func (a *Application) Start() error {
	if a.httpServer != nil {
		if err := a.httpServer.Start(); err != nil {
			return errors.Wrap(err, "http server start error")
		}
	}

	if a.grpcServer != nil {
		if err := a.grpcServer.Start(); err != nil {
			return errors.Wrap(err, "grpc server start error")
		}
	}

	return nil
}

func (a *Application) AwaitSignal() {
	c := make(chan os.Signal, 1)
	signal.Reset(syscall.SIGTERM, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	select {
	case s := <-c:
		a.logger.Info("receive a signal", zap.String("signal", s.String()))
		if a.httpServer != nil {
			if err := a.httpServer.Stop(); err != nil {
				a.logger.Warn("stop http server error", zap.Error(err))
			}
		}

		if a.grpcServer != nil {
			if err := a.grpcServer.Stop(); err != nil {
				a.logger.Warn("stop grpc server error", zap.Error(err))
			}
		}

		os.Exit(0)
	}
}
