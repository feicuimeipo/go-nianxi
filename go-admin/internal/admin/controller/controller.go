package controller

import (
	"gitee.com/go-nianxi/go-common/pkg/valid"
	"gitee.com/go-nianxi/go-admin/internal/admin/dao"
	"github.com/casbin/casbin/v2"
	"github.com/google/wire"
	"go.uber.org/zap"
)

type Controller struct {
	UserController         *UserController
	ApiController          *ApiController
	MenuController         *MenuController
	OperationLogController *OperationLogController
	RoleController         *RoleController
	ApplicationController  *ApplicationController
}

func New(dao *dao.Dao, valid *valid.Validator, casbinEnforcer *casbin.Enforcer, logger *zap.Logger) *Controller {
	var controller = new(Controller)
	controller.UserController = NewUserController(dao.UserDao, dao.RoleDao, valid)
	controller.ApiController = NewApiController(dao.ApiDao, dao.UserDao, dao.ApplicationDao, valid)
	controller.MenuController = NewMenuController(dao, valid, logger)
	controller.OperationLogController = NewOperationLogController(dao.OperationLogDao, valid)
	controller.RoleController = NewRoleController(dao, casbinEnforcer, valid)
	controller.ApplicationController = NewApplicationController(dao.ApplicationDao, valid)
	return controller
}

var ProviderSet = wire.NewSet(New)
