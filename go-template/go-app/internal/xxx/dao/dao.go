package dao

import (
	"github.com/google/wire"
	"gorm.io/gorm"
)

type Dao struct {
	HelloDao *HelloDao
}

func New(db *gorm.DB) *Dao {
	return &Dao{
		HelloDao: NewHelloDao(db),
	}
}

var ProviderSet = wire.NewSet(New)
