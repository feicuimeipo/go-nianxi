package dao

import (
	"errors"
	"fmt"
	"gitee.com/go-nianxi/go-admin/internal/admin/domain/vo"
	"gitee.com/go-nianxi/go-admin/internal/admin/model"
	"gorm.io/gorm"
	"strings"
)

type ApplicationDao struct {
	db *gorm.DB
}

func NewApplicationDao(db *gorm.DB) *ApplicationDao {
	return &ApplicationDao{
		db: db,
	}
}

func (a *ApplicationDao) GetAllApplicationTypes() ([]*model.ApplicationType, error) {
	var list []*model.ApplicationType
	err := a.db.Model(&model.ApplicationType{}).Order("created_at DESC").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (a *ApplicationDao) GetApplicationByTypeId(typeId uint) ([]*model.Application, error) {
	var list []*model.Application
	err := a.db.Model(&model.Application{}).Where("type_id=?", typeId).Find(&list).Order("created_at DESC").Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (a *ApplicationDao) GetAll() ([]*model.Application, error) {
	var list []*model.Application
	err := a.db.Model(&model.Application{}).Find(&list).Order("created_at DESC").Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (a *ApplicationDao) GetApplications(req *vo.ApplicationListRequest) ([]*model.Application, int64, error) {
	var list []*model.Application
	db := a.db.Model(&model.Application{}).Order("created_at DESC")

	appName := strings.TrimSpace(req.AppName)
	if appName != "" {
		db = db.Where("app_name LIKE ?", fmt.Sprintf("%s%%", appName))
	}
	alias := strings.TrimSpace(req.Alias)
	if alias != "" {
		db = db.Where("alias LIKE ?", fmt.Sprintf("%s%%", alias))
	}

	title := strings.TrimSpace(req.Title)
	if alias != "" {
		db = db.Where("title LIKE ?", fmt.Sprintf("%s%%", title))
	}

	baseUrl := strings.TrimSpace(req.BaseUrl)
	if alias != "" {
		db = db.Where("base_url LIKE ?", fmt.Sprintf("%s%%", baseUrl))
	}

	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := int(req.PageNum)
	pageSize := int(req.PageSize)
	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	return list, total, nil
}

func (a *ApplicationDao) GetApplicationById(id uint) (*model.Application, error) {
	var application *model.Application
	err := a.db.Where("id =? ", id).First(&application).Error
	if err != nil {
		return nil, err
	}
	return application, nil
}

// 创建用户
func (a *ApplicationDao) CreateApplication(application *model.Application) error {
	err := a.db.Create(application).Error
	return err
}

// 更新用户
func (a *ApplicationDao) UpdateApplicationById(id uint, application *model.Application) error {
	err := a.db.Model(application).Where("id=?", id).Updates(application).Error
	if err != nil {
		return err
	}

	return err
}

// 批量删除
func (a *ApplicationDao) BatchDeleteIds(ids []uint) error {
	// 用户和角色存在多对多关联关系
	var list []*model.Application
	err := a.db.Model(&model.Api{}).Where(" application_id in (?)", ids).Find(&list).Error
	if err != nil || len(list) == 0 {
		a.db.Model(&model.Application{}).Delete(" id in (?)", ids)
	} else {
		return errors.New("记录已在使用，不能删除")
	}
	return err
}
