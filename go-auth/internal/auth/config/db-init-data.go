package config

import (
	"errors"
	"gitee.com/go-nianxi/go-auth/internal/auth/model"
	"gitee.com/go-nianxi/go-common/pkg/base"
	"gitee.com/go-nianxi/go-common/pkg/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 初始化mysql数据
func initUserData(db *gorm.DB, logger *zap.SugaredLogger) {
	// 是否初始化数据
	defaultAvatar := "/default_avator.png"
	Introduction := ""
	users := []model.User{
		{
			Model:        gorm.Model{ID: 1},
			Username:     "admin",
			Password:     utils.GenPasswd("123456"),
			Mobile:       "18811111111",
			Email:        "auth_admin@abc.com",
			Avatar:       defaultAvatar,
			Introduction: Introduction,
			Status:       1,
			WxOpenId:     "",
		},
	}

	newUsers := make([]model.User, 0)
	for _, user := range users {
		user.Nickname = user.Username
		err := db.Where("username = ?", user.Username).First(&user).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newUsers = append(newUsers, user)
		}
	}

	if len(newUsers) > 0 {
		err := db.Create(&newUsers).Error
		if err != nil {
			logger.Errorf("写入用户数据失败：%v", err)
		}

	}

}

func InitData(db *gorm.DB, logger *zap.Logger) (bool, error) {
	if base.Conf.O.InitData {
		initUserData(db, logger.Sugar())
	}

	return true, nil
}

func dbAutoMigrate(db *gorm.DB) {
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&model.LoginAuditLog{},
		&model.User{})
}
