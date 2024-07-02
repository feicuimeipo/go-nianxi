package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username     string  `gorm:"type:varchar(20);not null;unique" json:"username"`
	Password     string  `gorm:"size:255;not null" json:"password"`
	Mobile       string  `gorm:"type:varchar(11);not null;unique" json:"mobile"`
	Avatar       string  `gorm:"type:varchar(255)" json:"avatar"`
	Nickname     *string `gorm:"type:varchar(20)" json:"nickname"`
	Introduction *string `gorm:"type:varchar(255)" json:"introduction"`
	Status       uint    `gorm:"type:tinyint(1);default:1;comment:'1正常, 2禁用'" json:"status"`
	Creator      string  `gorm:"type:varchar(20);" json:"creator"`
	Email        string  `gorm:"type:varchar(100);not null;default:'';unique;comment:'邮箱'"`
	RealName     string  `gorm:"type:varchar(20);not null;default:'';comment:'真名'"`
	IDCard       string  `gorm:"type:varchar(20);not null;default:'';comment:'身份证号'"`
	IsAdmin      bool    `gorm:"type:tinyint(1);default:0; comment:'0否 1 启动'"`
	WxOpenId     string  `gorm:"type:varchar(60);comment:'微信openId'" `
	TenantId     uint    `gorm:"type:varchar(100);not null;default:1;comment:'当前活跃的租户'"`
	Roles        []*Role `gorm:"many2many:sys_user_roles;" json:"roles"`
}
