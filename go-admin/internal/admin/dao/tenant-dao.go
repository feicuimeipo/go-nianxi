package dao

import (
	"gitee.com/go-nianxi/go-admin/internal/admin/model"
	"gorm.io/gorm"
)

type TenantDao struct {
	db *gorm.DB
}

func NewTenantDao(db *gorm.DB) *TenantDao {
	return &TenantDao{
		db: db,
	}
}

func (dao *TenantDao) GetTenantByTenantId(tenantId uint) (*model.TenantInfo, error) {
	var tenants model.TenantInfo
	err := dao.db.Where("id = ?", tenantId).Find(&tenants).Error
	return &tenants, err
}

func (dao *TenantDao) GetTenantsByUserId(userId uint) (*model.TenantInfo, error) {
	var model *model.TenantInfo
	err := dao.db.Joins("inner join user_tenant t on user_tenants.user_id = ? and user_tenants.tenant_id = tenant_info.id", userId).First(&model).Error //.Scan(&tenant)
	return model, err
}
