package model

import (
	"gorm.io/gorm"
	"time"
)

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
	WxOpenId     string  `gorm:"type:varchar(60);comment:'微信openId'" `
	TenantId     uint    `gorm:"type:varchar(100);not null;default:1;comment:'当前活跃的租户'"`
	Roles        []*Role `gorm:"many2many:sys_user_roles;" json:"roles"`
}

type Role struct {
	gorm.Model
	Name          string  `gorm:"type:varchar(20);not null;unique_index:ROLE_N_A" json:"name"`
	Keyword       string  `gorm:"type:varchar(20);not null;unique_index:ROLE_K_A;" json:"keyword"`
	Desc          *string `gorm:"type:varchar(100);"                 json:"desc"`
	Status        uint    `gorm:"type:tinyint(1);default:1;comment:'1正常, 2禁用'" json:"status"`
	Sort          uint    `gorm:"type:int(3);default:999;comment:'角色排序(排序越大权限越低, 不能查看比自己序号小的角色, 不能编辑同序号用户权限, 排序为1表示超级管理员)'" json:"sort"`
	Creator       string  `gorm:"type:varchar(20);" json:"creator"`
	ApplicationId uint    `gorm:"default:0;comment:应用编号;unique_index:ROLE_N_A;unique_index:ROLE_K_A" json:"applicationId"`
}

type OperationLog struct {
	gorm.Model
	Username   string    `gorm:"type:varchar(20);comment:'用户登录名'" json:"username"`
	Ip         string    `gorm:"type:varchar(20);comment:'Ip地址'" json:"ip"`
	IpLocation string    `gorm:"type:varchar(20);comment:'Ip所在地'" json:"ipLocation"`
	Method     string    `gorm:"type:varchar(20);comment:'请求方式'" json:"method"`
	Path       string    `gorm:"type:varchar(100);comment:'访问路径'" json:"path"`
	Desc       string    `gorm:"type:varchar(100);comment:'说明'" json:"desc"`
	Status     int       `gorm:"type:int(4);comment:'响应状态码'" json:"status"`
	StartTime  time.Time `gorm:"type:datetime(3);comment:'发起时间'" json:"startTime"`
	TimeCost   int64     `gorm:"type:int(6);comment:'请求耗时(ms)'" json:"timeCost"`
	UserAgent  string    `gorm:"type:varchar(20);comment:'浏览器标识'" json:"userAgent"`
}
