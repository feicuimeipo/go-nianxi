package model

import (
	"gorm.io/gorm"
)

type Hello struct {
	gorm.Model
	Msg string `gorm:"type:varchar(20);comment:'用户登录名'" json:"msg"`
}
