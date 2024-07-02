package dao

import (
	"gitee.com/go-nianxi/go-common/pkg/valid"
	"gitee.com/go-nianxi/go-template/internal/xxx/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type HelloDao struct {
	db     *gorm.DB
	valid  *valid.Validator
	logger *zap.Logger
}

func NewHelloDao(db *gorm.DB) *HelloDao {
	return &HelloDao{
		db: db,
	}
}

func (dao *HelloDao) GetHelloById(id uint) (*model.Hello, error) {
	return &model.Hello{
		Model: gorm.Model{ID: 1},
		Msg:   "hello",
	}, nil

}
