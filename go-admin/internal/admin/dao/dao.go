package dao

import (
	"github.com/casbin/casbin/v2"
	"github.com/google/wire"
	"gorm.io/gorm"
)

type Dao struct {
	RoleDao         *RoleDao
	UserDao         *UserDao
	ApiDao          *ApiDao
	TenantDao       *TenantDao
	MenuDao         *MenuDao
	OperationLogDao *OperationLogDao
	ApplicationDao  *ApplicationDao
}

func New(db *gorm.DB, casbinEnforcer *casbin.Enforcer) *Dao {
	return &Dao{
		RoleDao:         NewRoleDao(db, casbinEnforcer),
		UserDao:         NewUserDao(db),
		ApiDao:          NewApiDao(db, casbinEnforcer),
		TenantDao:       NewTenantDao(db),
		MenuDao:         NewMenuDao(db),
		OperationLogDao: NewOperationLogDao(db),
		ApplicationDao:  NewApplicationDao(db),
	}
}

var ProviderSet = wire.NewSet(New)
