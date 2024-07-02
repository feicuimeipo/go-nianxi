package repository

import (
	"errors"
	"fmt"
	"gitee.com/go-nianxi/go-common/pkg/cache"
	"gitee.com/go-nianxi/go-template/internal/admin/config"
	"gitee.com/go-nianxi/go-template/internal/admin/model"
	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
)

type SystemRepository struct {
	db    *gorm.DB
	cache cache.ICache
}

// SystemRepository构造函数
func NewSystemRepository() *SystemRepository {

	return &SystemRepository{
		db:    config.GetBaseDb(),
		cache: cache.NewLocalCacheService(),
	}
}

// 获取当前登录用户信息
// 需要缓存，减少数据库访问
func (ur *SystemRepository) GetCurrentUser(c *gin.Context) (*model.User, error) {
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
func (ur *SystemRepository) GetCurrentUserMinRoleSort(c *gin.Context) (uint, *model.User, error) {
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
func (ur *SystemRepository) GetUserById(id uint) (model.User, error) {
	fmt.Println("GetUserById---")
	var user model.User
	err := ur.db.Where("id = ?", id).Preload("Roles").First(&user).Error
	return user, err
}

// 根据用户ID获取用户角色排序最小值
func (ur *SystemRepository) GetUserMinRoleSortsByIds(ids []uint) ([]int, error) {
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
func (ur *SystemRepository) SetCache(username string, user model.User) {
	ur.cache.Set(username, user, cache.DefaultExpiration)
}

// 清理所有用户信息缓存
func (ur *SystemRepository) ClearCache() {
	ur.cache.GetLocalCache().Clear()
}

// var Logs []model.OperationLog //全局变量多个线程需要加锁，所以每个线程自己维护一个
// 处理OperationLogChan将日志记录到数据库
func (o SystemRepository) SaveOperationLogChannel(olc <-chan *model.OperationLog) {
	// 只会在线程开启的时候执行一次
	Logs := make([]model.OperationLog, 0)

	// 一直执行--收到olc就会执行
	for log := range olc {
		Logs = append(Logs, *log)
		// 每10条记录到数据库
		if len(Logs) > 5 {
			o.db.Create(&Logs)
			Logs = make([]model.OperationLog, 0)
		}
	}
}
