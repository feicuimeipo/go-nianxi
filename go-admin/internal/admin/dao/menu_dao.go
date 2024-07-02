package dao

import (
	"errors"
	"gitee.com/go-nianxi/go-admin/internal/admin/model"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
	"strings"
)

type MenuDao struct {
	db *gorm.DB
}

func NewMenuDao(db *gorm.DB) *MenuDao {
	return &MenuDao{
		db: db,
	}
}

// 获取菜单列表
func (m *MenuDao) GetMenus(applicationId uint) ([]*model.Menu, error) {
	var menus []*model.Menu
	err := m.db.Where("application_id =? ", applicationId).Order("sort").Find(&menus).Error
	return menus, err
}

// 获取菜单树
func (m *MenuDao) GetMenuTree(applicationId uint) ([]*model.Menu, error) {
	var menus []*model.Menu
	err := m.db.Where("application_id =? ", applicationId).Order("sort").Find(&menus).Error
	// parentId为0的是根菜单
	return GenMenuTree(0, menus, ""), err
}

func GenMenuTree(parentId uint, menus []*model.Menu, basePath string) []*model.Menu {
	tree := make([]*model.Menu, 0)

	for _, m := range menus {
		if *m.ParentId == parentId {
			m.BasePath = basePath
			childBasePath := ""
			if (strings.HasPrefix(m.Path, "http://") ||
				strings.HasPrefix(m.Path, "https://") ||
				strings.HasPrefix(m.Path, "mailto") ||
				strings.HasPrefix(m.Path, "tel:")) && basePath != "" {
				childBasePath = basePath
			} else {
				if basePath == "" {
					childBasePath = strings.TrimPrefix(m.Path, "/")
				} else {
					childBasePath = strings.TrimSuffix(basePath, "/") + "/" + strings.TrimPrefix(m.Path, "/")
				}
			}
			children := GenMenuTree(m.ID, menus, childBasePath)
			m.Children = children
			tree = append(tree, m)
		}
	}
	return tree
}

// 创建菜单
func (m *MenuDao) CreateMenu(menu *model.Menu) error {
	if menu.ApplicationId == 0 {
		return errors.New("应用不可以为空！")
	}
	err := m.db.Create(menu).Error
	return err
}

// 更新菜单
func (m *MenuDao) UpdateMenuById(menuId uint, menu *model.Menu) error {
	var old *model.Menu
	err := m.db.First(&old, menuId).Error
	if err != nil {
		return err
	}
	menu.ApplicationId = old.ApplicationId

	err = m.db.Model(menu).Where("id = ?", menuId).Updates(menu).Error
	return err
}

// 批量删除菜单
func (m *MenuDao) BatchDeleteMenuByIds(menuIds []uint) error {
	var menus []*model.Menu
	err := m.db.Where("id IN (?)", menuIds).Find(&menus).Error
	if err != nil {
		return err
	}
	err = m.db.Select("Roles").Unscoped().Delete(&menus).Error
	return err
}

// 根据用户ID获取用户的权限(可访问)菜单列表
func (m *MenuDao) GetUserMenusByUserId(userId uint, appId uint) ([]*model.Menu, error) {
	// 获取用户
	var user model.User
	err := m.db.Where("id = ?", userId).Preload("Roles").First(&user).Error
	if err != nil {
		return nil, err
	}
	// 获取角色
	roles := user.Roles
	// 所有角色的菜单集合
	allRoleMenus := make([]*model.Menu, 0)
	for _, role := range roles {
		var userRole model.Role
		if appId <= 0 {
			err = m.db.Where("id = ?", role.ID).Preload("Menus").Preload("Application").Order("application_id").First(&userRole).Error
			if err != nil {
				return nil, err
			}
		} else {
			err = m.db.Where("id = ? and application_id = ?", role.ID, appId).Preload("Menus").Preload("Application").First(&userRole).Error
			if err != nil {
				return nil, err
			}
		}
		// 获取角色的菜单
		menus := userRole.Menus
		allRoleMenus = append(allRoleMenus, menus...)
	}

	// 所有角色的菜单集合去重
	allRoleMenusId := make([]int, 0)
	for _, menu := range allRoleMenus {
		allRoleMenusId = append(allRoleMenusId, int(menu.ID))
	}
	allRoleMenusIdUniq := funk.UniqInt(allRoleMenusId)
	allRoleMenusUniq := make([]*model.Menu, 0)
	for _, id := range allRoleMenusIdUniq {
		for _, menu := range allRoleMenus {
			if id == int(menu.ID) {
				allRoleMenusUniq = append(allRoleMenusUniq, menu)
				break
			}
		}
	}

	// 获取状态status为1的菜单
	accessMenus := make([]*model.Menu, 0)
	for _, menu := range allRoleMenusUniq {
		if menu.Status == 1 {
			accessMenus = append(accessMenus, menu)
		}
	}

	return accessMenus, err
}

// 根据用户ID获取用户的权限(可访问)菜单树
func (m *MenuDao) GetUserMenuTreeByUserId(userId uint, appId uint) ([]*model.Menu, error) {
	menus, err := m.GetUserMenusByUserId(userId, appId)
	if err != nil {
		return nil, err
	}
	tree := GenMenuTree(0, menus, "")
	return tree, err
}
