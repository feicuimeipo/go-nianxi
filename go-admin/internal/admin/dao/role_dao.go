package dao

import (
	"errors"
	"fmt"
	"gitee.com/go-nianxi/go-admin/internal/admin/domain/vo"
	"gitee.com/go-nianxi/go-admin/internal/admin/model"
	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"

	"strings"
)

type RoleDao struct {
	db     *gorm.DB
	casbin *casbin.Enforcer
}

func NewRoleDao(db *gorm.DB, casbinEnforcer *casbin.Enforcer) *RoleDao {
	return &RoleDao{
		db:     db,
		casbin: casbinEnforcer,
	}
}

// 获取角色列表
func (r RoleDao) GetRoles(req *vo.RoleListRequest) ([]model.Role, int64, error) {
	var list []model.Role
	db := r.db.Model(&model.Role{}).Order("created_at DESC")

	if req.ApplicationId > 0 {
		db = db.Where("application_id = ?", req.ApplicationId)
	}

	name := strings.TrimSpace(req.Name)
	if name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}

	keyword := strings.TrimSpace(req.Keyword)
	if keyword != "" {
		db = db.Where("keyword LIKE ?", fmt.Sprintf("%%%s%%", keyword))
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
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	return list, total, err
}

// 根据角色ID获取角色
func (r RoleDao) GetRolesById(roleId uint) (*model.Role, error) {
	var role model.Role
	err := r.db.Where("id = ", roleId).Find(&role).Error
	return &role, err
}

// 根据角色ID获取角色
func (r RoleDao) GetRolesByIds(roleIds []uint) ([]*model.Role, error) {
	var list []*model.Role
	err := r.db.Where("id IN (?)", roleIds).Find(&list).Error
	return list, err
}

// 创建角色
func (r RoleDao) CreateRole(role *model.Role) error {

	err := r.db.Create(role).Error
	return err
}

// 更新角色
func (r RoleDao) UpdateRoleById(roleId uint, role *model.Role) error {
	var old *model.Menu
	err := r.db.First(&old, roleId).Error
	if err != nil {
		return err
	}
	role.ApplicationId = old.ApplicationId

	err = r.db.Model(&model.Role{}).Where("id = ?", roleId).Updates(role).Error
	return err
}

// 获取角色的权限菜单
func (r RoleDao) GetRoleMenusById(roleId uint) ([]*model.Menu, error) {
	var role model.Role
	err := r.db.Where("id = ?", roleId).Preload("Menus").First(&role).Error
	return role.Menus, err
}

// 更新角色的权限菜单
func (r RoleDao) UpdateRoleMenus(role *model.Role) error {
	err := r.db.Model(role).Association("Menus").Replace(role.Menus)
	return err
}

// 根据角色关键字获取角色的权限接口
func (r RoleDao) GetRoleApisByRoleKeyword(roleKeyword string) ([]*model.Api, error) {
	policies := r.casbin.GetFilteredPolicy(0, roleKeyword)

	// 获取所有接口
	var apis []*model.Api
	err := r.db.Find(&apis).Error
	if err != nil {
		return apis, errors.New("获取角色的权限接口失败")
	}

	accessApis := make([]*model.Api, 0)

	for _, policy := range policies {
		path := policy[1]
		method := policy[2]
		for _, api := range apis {
			if path == api.Path && method == api.Method {
				accessApis = append(accessApis, api)
				break
			}
		}
	}

	return accessApis, err

}

// 更新角色的权限接口（先全部删除再新增）
func (r RoleDao) UpdateRoleApis(roleKeyword string, reqRolePolicies [][]string) error {
	// 先获取path中的角色ID对应角色已有的police(需要先删除的)
	err := r.casbin.LoadPolicy()
	if err != nil {
		return errors.New("角色的权限接口策略加载失败")
	}
	rmPolicies := r.casbin.GetFilteredPolicy(0, roleKeyword)
	if len(rmPolicies) > 0 {
		isRemoved, _ := r.casbin.RemovePolicies(rmPolicies)
		if !isRemoved {
			return errors.New("更新角色的权限接口失败")
		}
	}
	isAdded, _ := r.casbin.AddPolicies(reqRolePolicies)
	if !isAdded {
		return errors.New("更新角色的权限接口失败")
	}
	err = r.casbin.LoadPolicy()
	if err != nil {
		return errors.New("更新角色的权限接口成功，角色的权限接口策略加载失败")
	} else {
		return err
	}
}

// 删除角色
func (r RoleDao) BatchDeleteRoleByIds(roleIds []uint) error {
	var roles []*model.Role
	err := r.db.Where("id IN (?)", roleIds).Find(&roles).Error
	if err != nil {
		return err
	}
	err = r.db.Select("Users", "Menus").Unscoped().Delete(&roles).Error
	// 删除成功就删除casbin policy
	if err == nil {
		for _, role := range roles {
			roleKeyword := role.Keyword
			rmPolicies := r.casbin.GetFilteredPolicy(0, roleKeyword)
			if len(rmPolicies) > 0 {
				isRemoved, _ := r.casbin.RemovePolicies(rmPolicies)
				if !isRemoved {
					return errors.New("删除角色成功, 删除角色关联权限接口失败")
				}
			}
		}

	}
	return err
}
