package repository

import (
	"errors"
	"fmt"
	"gitee.com/go-nianxi/go-template/internal/admin/domain/vo"
	"gitee.com/go-nianxi/go-template/internal/admin/model"
	"gorm.io/gorm"
	"strings"
)

type HelloRepository struct {
	db *gorm.DB
}

func NewHelloRepository(db *gorm.DB) *HelloRepository {
	return &HelloRepository{
		db: db,
	}
}

func (o HelloRepository) GetHellos(req *vo.HelloListRequest) ([]*model.Hello, int64, error) {

	var list []*model.Hello
	db := o.db.Model(&model.OperationLog{}).Order("start_time DESC")

	msg := strings.TrimSpace(req.Msg)
	if msg != "" {
		db = db.Where("msg LIKE ?", fmt.Sprintf("%%%s%%", msg))
	}

	// 分页
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := req.PageNum
	pageSize := req.PageSize
	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}

	return list, total, err

	return list, 1, nil

}

// 根据接口ID获取接口列表
func (a HelloRepository) GetApisById(apiIds []uint) ([]*model.Hello, error) {
	var apis []*model.Hello
	err := a.db.Where("id IN (?)", apiIds).Find(&apis).Error
	return apis, err
}

// 创建接口
func (a HelloRepository) CreateApi(api *model.Hello) error {

	err := a.db.Create(api).Error
	return err
}

// 更新接口
func (a HelloRepository) UpdateApiById(id uint, api *model.Hello) error {

	err := a.db.Model(api).Where("id = ?", id).Updates(api).Error
	if err != nil {
		return err
	}

	return err
}

// 批量删除接口
func (a HelloRepository) BatchDeleteApiByIds(apiIds []uint) error {

	apis, err := a.GetApisById(apiIds)
	if err != nil {
		return errors.New("根据接口ID获取接口列表失败")
	}
	if len(apis) == 0 {
		return errors.New("根据接口ID未获取到接口列表")
	}

	err = a.db.Where("id IN (?)", apiIds).Unscoped().Delete(&model.Hello{}).Error

	return err
}

// 根据接口路径和请求方式获取接口描述
func (a HelloRepository) GetApiDescByPath(path string, method string) (*model.Hello, error) {
	var hello model.Hello
	err := a.db.Where("path = ?", path).Where("method = ?", method).First(&hello).Error
	return &hello, err
}
