package controller

import (
	"gitee.com/go-nianxi/go-common/pkg/valid"
	"gitee.com/go-nianxi/go-template/internal/admin/repository"
	"github.com/google/wire"
	"go.uber.org/zap"
)

type Controller struct {
	HelloController *HelloController
}

func New(dao *repository.Repository, valid *valid.Validator, logger *zap.Logger) *Controller {
	var controller = new(Controller)

	controller.HelloController = NewHelloController(dao.HelloRepository, valid)

	return controller
}

var ProviderSet = wire.NewSet(New)
