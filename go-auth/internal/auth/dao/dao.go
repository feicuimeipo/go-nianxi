package dao

import (
	"github.com/google/wire"
	"gorm.io/gorm"
)

type Dao struct {
	UserDao          *UserDao
	LoginAuditLogDao *LoginAuditLogDao
}

func New(db *gorm.DB) *Dao {
	return &Dao{
		UserDao:          NewUserDao(db),
		LoginAuditLogDao: NewLoginAuditLogDao(db),
	}
}

type WriterOperationStatus struct {
	status bool
}

var ProviderSet = wire.NewSet(New)
