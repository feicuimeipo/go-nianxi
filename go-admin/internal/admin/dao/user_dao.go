package dao

import (
	"errors"
	"fmt"
	"gitee.com/go-nianxi/go-common/pkg/cache"
	"gitee.com/go-nianxi/go-admin/internal/admin/domain/vo"
	"gitee.com/go-nianxi/go-admin/internal/admin/model"
	"gitee.com/go-nianxi/go-admin/internal/admin/util"
	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
	"strings"
)

type UserDao struct {
	db    *gorm.DB
	cache cache.ICache
}

// UserDao构造函数
func NewUserDao(db *gorm.DB) *UserDao {

	return &UserDao{
		db:    db,
		cache: cache.NewLocalCacheService(),
	}
}

// 登录
func (ur *UserDao) Login(reqUser *model.User) (*model.User, error) {
	// 根据用户名获取用户(正常状态:用户状态正常)
	var user model.User
	err := ur.db.
		Where("username = ?", reqUser.Username).
		Preload("Roles").
		First(&user).Error
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 校验密码
	err = util.ComparePasswd(user.Password, reqUser.Password)
	if err != nil {
		return &user, errors.New("密码错误")
	}

	// 判断用户的状态
	userStatus := user.Status
	if userStatus != 1 {
		return nil, errors.New("用户被禁用")
	}

	// 判断用户拥有的所有角色的状态,全部角色都被禁用则不能登录
	roles := user.Roles
	isValidate := false
	for _, role := range roles {
		// 有一个正常状态的角色就可以登录
		if role.Status == 1 {
			isValidate = true
			break
		}
	}

	if !isValidate {
		return nil, errors.New("用户角色被禁用")
	}

	return &user, nil
}

// 获取当前登录用户信息
// 需要缓存，减少数据库访问
func (ur *UserDao) GetCurrentUser(c *gin.Context) (*model.User, error) {
	//var newUser model.User
	ctxUser, exist := c.Get("user")
	if !exist {
		return nil, errors.New("用户未登录")
	}
	u, _ := ctxUser.(model.User)

	// 先获取缓存
	cacheUser, found := ur.cache.Get(u.Username)
	var user model.User
	var err error
	if found {
		user = cacheUser.(model.User)
		err = nil
	} else {
		// 缓存中没有就获取数据库
		user, err = ur.GetUserById(u.ID)
		// 获取成功就缓存
		if err != nil {
			ur.cache.Delete(u.Username)
		} else {
			ur.cache.Set(u.Username, user, cache.DefaultExpiration)
		}
	}
	return &user, err
}

// 获取当前用户角色排序最小值（最高等级角色）以及当前用户信息
func (ur *UserDao) GetCurrentUserMinRoleSort(c *gin.Context) (uint, *model.User, error) {
	// 获取当前用户
	ctxUser, err := ur.GetCurrentUser(c)
	if err != nil {
		return 999, ctxUser, err
	}
	// 获取当前用户的所有角色
	currentRoles := ctxUser.Roles
	// 获取当前用户角色的排序，和前端传来的角色排序做比较
	var currentRoleSorts []int
	for _, role := range currentRoles {
		currentRoleSorts = append(currentRoleSorts, int(role.Sort))
	}
	// 当前用户角色排序最小值（最高等级角色）
	currentRoleSortMin := uint(funk.MinInt(currentRoleSorts))

	return currentRoleSortMin, ctxUser, nil
}

// 获取单个用户
func (ur *UserDao) GetUserById(id uint) (model.User, error) {
	fmt.Println("GetUserById---")
	var user model.User
	err := ur.db.Where("id = ?", id).Preload("Roles").First(&user).Error
	return user, err
}

// 获取用户列表
func (ur *UserDao) GetUsers(req *vo.UserListRequest) ([]*model.User, int64, error) {
	var list []*model.User
	db := ur.db.Model(&model.User{}).Order("created_at DESC")

	username := strings.TrimSpace(req.Username)
	if username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	nickname := strings.TrimSpace(req.Nickname)
	if nickname != "" {
		db = db.Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", nickname))
	}
	mobile := strings.TrimSpace(req.Mobile)
	if mobile != "" {
		db = db.Where("mobile LIKE ?", fmt.Sprintf("%%%s%%", mobile))
	}
	status := req.Status
	if status != 0 {
		db = db.Where("status = ?", status)
	}
	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := int(req.PageNum)
	pageSize := int(req.PageSize)
	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Preload("Roles").Find(&list).Error
	} else {
		err = db.Preload("Roles").Find(&list).Error
	}
	return list, total, err
}

// 更新密码
func (ur *UserDao) ChangePwd(username string, hashNewPasswd string) error {
	err := ur.db.Model(&model.User{}).Where("username = ?", username).Update("password", hashNewPasswd).Error
	// 如果更新密码成功，则更新当前用户信息缓存
	// 先获取缓存
	cacheUser, found := ur.cache.Get(username)
	if err == nil {
		if found {
			user := cacheUser.(model.User)
			user.Password = hashNewPasswd
			ur.cache.Set(username, user, cache.DefaultExpiration)
		} else {
			// 没有缓存就获取用户信息缓存
			var user model.User
			ur.db.Where("username = ?", username).First(&user)
			ur.cache.Set(username, user, cache.DefaultExpiration)
		}
	}

	return err
}

// 创建用户
func (ur *UserDao) CreateUser(user *model.User) error {
	err := ur.db.Create(user).Error
	return err
}

// 更新用户
func (ur *UserDao) UpdateUser(user *model.User) error {
	err := ur.db.Model(user).Updates(user).Error
	if err != nil {
		return err
	}
	err = ur.db.Model(user).Association("Roles").Replace(user.Roles)

	//err := ur.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&user).Error

	// 如果更新成功就更新用户信息缓存
	if err == nil {
		ur.cache.Set(user.Username, *user, cache.DefaultExpiration)
	}
	return err
}

// 批量删除
func (ur *UserDao) BatchDeleteUserByIds(ids []uint) error {
	// 用户和角色存在多对多关联关系
	var users []model.User
	for _, id := range ids {
		// 根据ID获取用户
		user, err := ur.GetUserById(id)
		if err != nil {
			return errors.New(fmt.Sprintf("未获取到ID为%d的用户", id))
		}
		users = append(users, user)
	}

	err := ur.db.Select("Roles").Unscoped().Delete(&users).Error
	// 删除用户成功，则删除用户信息缓存
	if err == nil {
		for _, user := range users {
			ur.cache.Delete(user.Username)
		}
	}
	return err
}

// 根据用户ID获取用户角色排序最小值
func (ur *UserDao) GetUserMinRoleSortsByIds(ids []uint) ([]int, error) {
	// 根据用户ID获取用户信息
	var userList []model.User
	err := ur.db.Where("id IN (?)", ids).Preload("Roles").Find(&userList).Error
	if err != nil {
		return []int{}, err
	}
	if len(userList) == 0 {
		return []int{}, errors.New("未获取到任何用户信息")
	}
	var roleMinSortList []int
	for _, user := range userList {
		roles := user.Roles
		var roleSortList []int
		for _, role := range roles {
			roleSortList = append(roleSortList, int(role.Sort))
		}
		roleMinSort := funk.MinInt(roleSortList)
		roleMinSortList = append(roleMinSortList, roleMinSort)
	}
	return roleMinSortList, nil
}

// 设置用户信息缓存
func (ur *UserDao) SetCache(username string, user model.User) {
	ur.cache.Set(username, user, cache.DefaultExpiration)
}

// 根据角色ID更新拥有该角色的用户信息缓存
func (ur *UserDao) UpdateCacheByRoleId(roleId uint) error {

	var role model.Role
	err := ur.db.Where("id = ?", roleId).Preload("Users").First(&role).Error
	if err != nil {
		return errors.New("根据角色ID角色信息失败")
	}

	users := role.Users
	if len(users) == 0 {
		return errors.New("根据角色ID未获取到拥有该角色的用户")
	}

	// 更新用户信息缓存
	for _, user := range users {
		_, found := ur.cache.Get(user.Username)
		if found {
			ur.cache.Set(user.Username, *user, cache.DefaultExpiration)
		}
	}

	return err
}

// 清理所有用户信息缓存
func (ur *UserDao) ClearCache() {
	ur.cache.GetLocalCache().Clear()
}
