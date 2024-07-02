package repository

import (
	"github.com/google/wire"
	"gorm.io/gorm"
)

type Repository struct {
	HelloRepository  *HelloRepository
	SystemRepository *SystemRepository
}

func New(db *gorm.DB) *Repository {
	return &Repository{
		HelloRepository:  NewHelloRepository(db),
		SystemRepository: NewSystemRepository(),
	}
}

var ProviderSet = wire.NewSet(New)
