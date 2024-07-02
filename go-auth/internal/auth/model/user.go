package model

import (
	auth_common "gitee.com/go-nianxi/go-auth/pkg/auth-client/auth-common"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `gorm:"type:varchar(20); null;comment:'用户名'"`
	Password     string `gorm:"size:255;not null;comment:'密码'"`
	Mobile       string `gorm:"type:varchar(11);not null;default;unique'';unique;comment:'手机号'"`
	Email        string `gorm:"type:varchar(100);null;default:'';comment:'邮箱'"`
	Avatar       string `gorm:"type:varchar(255);comment:'头像'" `
	Nickname     string `gorm:"type:varchar(20); null;default:'';comment:'昵称'"`
	Status       uint   `gorm:"type:tinyint(1);default:1; comment:'1正常, 2禁用'" `
	Introduction string `gorm:"type:varchar(255);comment:'描述'" `
	WxOpenId     string `gorm:"type:varchar(60);comment:'微信openId'" `
}

func ToAuthUser(u *User) *auth_common.AuthUser {
	dto := auth_common.AuthUser{}
	dto.ID = u.ID
	dto.Username = u.Username
	dto.Mobile = u.Mobile
	dto.Email = u.Email
	dto.Avatar = u.Avatar
	dto.WxOpenId = u.WxOpenId
	dto.Introduction = u.Introduction
	dto.Nickname = u.Nickname

	return &dto
}
