package model

import (
	"gorm.io/gorm"
)

type Hello struct {
	gorm.Model
	Msg string `gorm:"type:varchar(20);comment:'hello'" json:"hello"`
}
