package model

import "gorm.io/gorm"

type ApplicationType struct {
	gorm.Model
	Desc string `gorm:"type:varchar(64);  not null;unique;      comment:描述';"        json:"Desc"`
}
type Application struct {
	gorm.Model
	AppName      string  `gorm:"type:varchar(64);  not null;unique;      comment:应用名,对应config.yaml文件中的app.Name';"        json:"appName"`
	Alias        string  `gorm:"type:varchar(64);  not null;unique;      comment:应用名,对应config.yaml文件中的app.Name';"        json:"alias"`
	Title        string  `gorm:"type:varchar(30);  not null;unique;  	 comment:'中文-应用标题';"                                json:"Title"`
	BaseUrl      string  `gorm:"type:varchar(128); not null;unique;      comment:'管理端访问主url';"                             json:"baseUrl"`
	TypeId       uint    `gorm:"type:tinyint(1);  not null;default:0;    comment:'字典表中字段-ApplicationType.Id(0为基础应用)';"  json:"typeId"`
	Introduction *string `gorm:"type:varchar(255);                       comment:'描述'"                                        json:"introduction"`
	Icon         string  `gorm:"type:varchar(255);"                json:"Icon"`
}
