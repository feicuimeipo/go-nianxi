package dao

import (
	"fmt"
	"gitee.com/go-nianxi/go-auth/internal/auth/model"
	"gorm.io/gorm"
	"strings"
)

type UserDao struct {
	db *gorm.DB
}

const (
	AccountLoginType_Mobile = iota
	AccountLoginType_Email
	AccountLoginType_USER
)

// UserDao构造函数
func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (dao *UserDao) GetUserForLogin(loginName string, accountType int) (*model.User, error) {
	search := model.User{}
	if accountType == AccountLoginType_Mobile {
		search.Mobile = loginName
	} else if accountType == AccountLoginType_Email {
		search.Email = loginName
	} else {
		search.Username = loginName
	}
	var user = model.User{}
	err := dao.db.Where(search).First(&user).Error
	if err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

func (dao *UserDao) GetUsersByMobile(mobile string) ([]*model.User, error) {
	var users []*model.User
	search := model.User{Mobile: mobile}
	err := dao.db.Where(search).Find(users).Error
	if err != nil {
		return nil, err
	} else {
		return users, nil
	}
}

func (dao *UserDao) GetUsersByNickName(nickName string) ([]*model.User, error) {
	fmt.Println("GetUserByEmail---")
	var users []*model.User
	err := dao.db.Where("nick_name = ? and user_type", nickName).Find(&users).Error
	return users, err
}

// 获取单个用户
func (dao *UserDao) GetUserById(id uint) (*model.User, error) {
	fmt.Println("GetUserById---")
	var user model.User
	err := dao.db.Where("id = ?", id).First(&user).Error
	return &user, err
}

// 获取单个用户
func (dao *UserDao) GetUserByUserName(userName string) (*model.User, error) {
	fmt.Println("GetUserById---")
	var user model.User
	err := dao.db.Where("username = ? and user_type = ?", userName).First(&user).Error
	return &user, err
}

// 用户注册
func (dao *UserDao) CreateUser(user *model.User) (*model.User, error) {
	user.Status = 1
	err := dao.db.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// 更新用户
func (dao *UserDao) UpdateUser(user *model.User) error {
	err := dao.db.Model(user).Updates(user).Error
	if err != nil {
		return err
	}
	return err
}

// 获取用户列表
func (dao *UserDao) GetUserListLimit5(req *model.User) ([]*model.User, error) {
	var list []*model.User

	db := dao.db.Model(&model.User{}).Order("created_at DESC")

	username := strings.TrimSpace(req.Username)
	if username != "" {
		db = db.Where("username = ?", username)
	}
	nickname := strings.TrimSpace(req.Nickname)
	if nickname != "" {
		db = db.Where("nickname = ? ", nickname)
	}
	mobile := strings.TrimSpace(req.Mobile)
	if mobile != "" {
		db = db.Where("mobile = ?", mobile)
	}
	email := strings.TrimSpace(req.Email)
	if email != "" {
		db = db.Where("email = ? ", email)
	}

	err := db.Limit(5).Find(&list).Error
	return list, err
}
