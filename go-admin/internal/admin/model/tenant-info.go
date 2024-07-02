package model

import "gorm.io/gorm"

type TenantInfo struct {
	gorm.Model
	Name         string  `gorm:"type:varchar(20);unique;not null;comment:'请求方式'"   json:"name"`
	TenantCode   string  `gorm:"type:varchar(50);unique;not null;comment:'访问路径'"   json:"tenantCode"`
	Introduction *string `gorm:"type:varchar(255);comment:'描述'" json:"introduction"`
}
