package model

import (
	"gorm.io/gorm"
	"time"
)

type LoginAuditLog struct {
	gorm.Model
	Username   string    `gorm:"type:varchar(20);comment:'用户登录名'" json:"username"`
	Ip         string    `gorm:"type:varchar(20);comment:'Ip地址'"    json:"ip"`
	IpLocation string    `gorm:"type:varchar(20);comment:'Ip所在地'"  json:"ipLocation"`
	Method     string    `gorm:"type:varchar(20);comment:'请求方式'"   json:"method"`
	Path       string    `gorm:"type:varchar(100);comment:'访问路径'"  json:"path"`
	Status     int       `gorm:"type:int(4);comment:'响应状态码'"      json:"status"`
	StartTime  time.Time `gorm:"type:datetime(3);comment:'发起时间'"   json:"startTime"`
	TimeCost   int64     `gorm:"type:int(6);comment:'请求耗时(ms)'"    json:"timeCost"`
	UserAgent  string    `gorm:"type:varchar(20);comment:'浏览器标识'"  json:"userAgent"`
	Action     string    `gorm:"type:varchar(20);comment:'动作'"  json:"action"`
}
