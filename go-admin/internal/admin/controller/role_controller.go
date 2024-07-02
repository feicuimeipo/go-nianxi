package controller

import (
	"fmt"
	"gitee.com/go-nianxi/go-common/pkg/ecode/resp"
	"gitee.com/go-nianxi/go-common/pkg/valid"
	"gitee.com/go-nianxi/go-admin/internal/admin/dao"
	"gitee.com/go-nianxi/go-admin/internal/admin/domain/vo"
	"gitee.com/go-nianxi/go-admin/internal/admin/model"
	"gitee.com/go-nianxi/go-admin/internal/pkg/ecode"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
	"strconv"
)

type RoleController struct {
	roleDao        *dao.RoleDao
	userDao        *dao.UserDao
	menuDao        *dao.MenuDao
	apiDao         *dao.ApiDao
	appDao         *dao.ApplicationDao
	valid          *valid.Validator
	casbinEnforcer *casbin.Enforcer
}

func NewRoleController(dao *dao.Dao, casbinEnforcer *casbin.Enforcer, valid *valid.Validator) *RoleController {
	return &RoleController{
		userDao:        dao.UserDao,
		roleDao:        dao.RoleDao,
		menuDao:        dao.MenuDao,
		apiDao:         dao.ApiDao,
		appDao:         dao.ApplicationDao,
		valid:          valid,
		casbinEnforcer: casbinEnforcer,
	}
}

// 获取角色列表
func (rc RoleController) GetRoles(c *gin.Context) {
	var req vo.RoleListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, rc.valid.Translate(err))
		return
	}

	if req.ApplicationId == 0 {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "2.应用编号不可以为空")
		return
	}

	// 获取角色列表
	roles, total, err := rc.roleDao.GetRoles(&req)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "获取角色列表失败: "+err.Error())
		return
	}
	resp.Success(c, gin.H{"roles": roles, "total": total}, "获取角色列表成功")
}

// 获取角色列表
func (rc RoleController) GetAllRoles(c *gin.Context) {
	req := vo.RoleListRequest{}

	// 获取角色列表
	roles, total, err := rc.roleDao.GetRoles(&req)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "获取角色列表失败: "+err.Error())
		return
	}
	resp.Success(c, gin.H{"roles": roles, "total": total}, "获取角色列表成功")
}

// 创建角色
func (rc RoleController) CreateRole(c *gin.Context) {
	var req vo.CreateRoleRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, rc.valid.Translate(err))
		return
	}

	// 获取当前用户最高角色等级
	sort, ctxUser, err := rc.userDao.GetCurrentUserMinRoleSort(c)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "获取当前用户最高角色等级失败: "+err.Error())
		return
	}

	// 用户不能创建比自己等级高或相同等级的角色
	if sort >= req.Sort {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "不能创建比自己等级高或相同等级的角色")
		return
	}

	if req.ApplicationId == 0 {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, "3.应用编号不可以为空！")
		return
	}

	app, err := rc.appDao.GetApplicationById(req.ApplicationId)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, err.Error())
		return
	}

	role := model.Role{
		Name:          req.Name,
		Keyword:       req.Keyword,
		Desc:          &req.Desc,
		Status:        req.Status,
		Sort:          req.Sort,
		ApplicationId: app.ID,
		Creator:       ctxUser.Username,
	}

	// 创建角色
	err = rc.roleDao.CreateRole(&role)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "创建角色失败: "+err.Error())
		return
	}
	resp.OK(c, nil, "创建角色成功")

}

// 更新角色
func (rc RoleController) UpdateRoleById(c *gin.Context) {
	var req vo.CreateRoleRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, rc.valid.Translate(err))
		return
	}

	// 获取path中的roleId
	roleId, _ := strconv.Atoi(c.Param("roleId"))
	if roleId <= 0 {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "角色ID不正确")
		return
	}

	// 当前用户角色排序最小值（最高等级角色）以及当前用户

	minSort, ctxUser, err := rc.userDao.GetCurrentUserMinRoleSort(c)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, err.Error())
		return
	}

	// 不能更新比自己角色等级高或相等的角色
	// 根据path中的角色ID获取该角色信息
	roles, err := rc.roleDao.GetRolesByIds([]uint{uint(roleId)})
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, err.Error())
		return
	}
	if len(roles) == 0 {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "未获取到角色信息")
		return
	}
	if minSort >= roles[0].Sort {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "不能更新比自己角色等级高或相等的角色")
		return
	}

	// 不能把角色等级更新得比当前用户的等级高
	if minSort >= req.Sort {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "不能把角色等级更新得比当前用户的等级高或相同")
		return
	}

	role := model.Role{
		Name:          req.Name,
		Keyword:       req.Keyword,
		Desc:          &req.Desc,
		Status:        req.Status,
		Sort:          req.Sort,
		ApplicationId: req.ApplicationId,
		Creator:       ctxUser.Username,
	}

	// 更新角色
	err = rc.roleDao.UpdateRoleById(uint(roleId), &role)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "更新角色失败: "+err.Error())
		return
	}

	// 如果更新成功，且更新了角色的keyword, 则更新casbin中policy
	if req.Keyword != roles[0].Keyword {
		// 获取policy
		rolePolicies := rc.casbinEnforcer.GetFilteredPolicy(0, roles[0].Keyword)
		if len(rolePolicies) == 0 {
			resp.OK(c, nil, "更新角色成功")
			return
		}
		rolePoliciesCopy := make([][]string, 0)
		// 替换keyword
		for _, policy := range rolePolicies {
			policyCopy := make([]string, len(policy))
			copy(policyCopy, policy)
			rolePoliciesCopy = append(rolePoliciesCopy, policyCopy)
			policy[0] = req.Keyword
		}

		//gormadapter未实现UpdatePolicies方法，等gorm更新---
		//isUpdated, _ := common.CasbinEnforcer.UpdatePolicies(rolePoliciesCopy, rolePolicies)
		//if !isUpdated {
		//	resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "更新角色成功，但角色关键字关联的权限接口更新失败！")
		//	return
		//}

		// 这里需要先新增再删除（先删除再增加会出错）
		isAdded, _ := rc.casbinEnforcer.AddPolicies(rolePolicies)
		if !isAdded {
			resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "更新角色成功，但角色关键字关联的权限接口更新失败")
			return
		}
		isRemoved, _ := rc.casbinEnforcer.RemovePolicies(rolePoliciesCopy)
		if !isRemoved {
			resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "更新角色成功，但角色关键字关联的权限接口更新失败")
			return
		}
		err := rc.casbinEnforcer.LoadPolicy()
		if err != nil {
			resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "更新角色成功，但角色关键字关联角色的权限接口策略加载失败")
			return
		}

	}

	// 更新角色成功处理用户信息缓存有两种做法:（这里使用第二种方法，因为一个角色下用户数量可能很多，第二种方法可以分散数据库压力）
	// 1.可以帮助用户更新拥有该角色的用户信息缓存,使用下面方法
	// err = ur.UpdateUserInfoCacheByRoleId(uint(roleId))
	// 2.直接清理缓存，让活跃的用户自己重新缓存最新用户信息
	//ur.ClearUserInfoCache()
	rc.userDao.ClearCache()
	resp.OK(c, nil, "更新角色成功")
}

// 获取角色的权限菜单
func (rc RoleController) GetRoleMenusById(c *gin.Context) {
	// 获取path中的roleId
	roleId, _ := strconv.Atoi(c.Param("roleId"))
	if roleId <= 0 {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "角色ID不正确")
		return
	}
	menus, err := rc.roleDao.GetRoleMenusById(uint(roleId))
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "获取角色的权限菜单失败: "+err.Error())
		return
	}
	resp.Success(c, gin.H{"menus": menus}, "获取角色的权限菜单成功")
	//resp.Success(c, menus)
}

// 更新角色的权限菜单
func (rc RoleController) UpdateRoleMenusById(c *gin.Context) {
	var req vo.UpdateRoleMenusRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, rc.valid.Translate(err))
		return
	}

	// 获取path中的roleId
	roleId, _ := strconv.Atoi(c.Param("roleId"))
	if roleId <= 0 {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "角色ID不正确")
		return
	}
	//根据path中的角色ID获取该角色信息
	role, err := rc.roleDao.GetRolesById(uint(roleId))
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, err.Error())
		return
	}

	// 当前用户角色排序最小值（最高等级角色）以及当前用户

	minSort, ctxUser, err := rc.userDao.GetCurrentUserMinRoleSort(c)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, err.Error())
		return
	}

	// (非管理员)不能更新比自己角色等级高或相等角色的权限菜单
	if minSort != 1 {
		if minSort >= role.Sort {
			resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "不能更新比自己角色等级高或相等角色的权限菜单")
			return
		}
	}

	// 获取当前用户所拥有的权限菜单
	ctxUserMenus, err := rc.menuDao.GetUserMenusByUserId(ctxUser.ID, role.ApplicationId)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "获取当前用户的可访问菜单列表失败: "+err.Error())
		return
	}

	// 获取当前用户所拥有的权限菜单ID
	ctxUserMenusIds := make([]uint, 0)
	for _, menu := range ctxUserMenus {
		ctxUserMenusIds = append(ctxUserMenusIds, menu.ID)
	}

	// 前端传来最新的MenuIds集合
	menuIds := req.MenuIds

	// 用户需要修改的菜单集合
	reqMenus := make([]*model.Menu, 0)

	// (非管理员)不能把角色的权限菜单设置的比当前用户所拥有的权限菜单多
	if minSort != 1 {
		for _, id := range menuIds {
			if !funk.Contains(ctxUserMenusIds, id) {
				resp.Fail(c, ecode.RECODE_INTERNAL_ERR, fmt.Sprintf("无权设置ID为%d的菜单", id))
				return
			}
		}

		for _, id := range menuIds {
			for _, menu := range ctxUserMenus {
				if id == menu.ID {
					reqMenus = append(reqMenus, menu)
					break
				}
			}
		}
	} else {
		// 管理员随意设置
		// 根据menuIds查询查询菜单
		menus, err := rc.menuDao.GetMenus(role.ApplicationId)
		if err != nil {
			resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "获取菜单列表失败: "+err.Error())
			return
		}
		for _, menuId := range menuIds {
			for _, menu := range menus {
				if menuId == menu.ID {
					reqMenus = append(reqMenus, menu)
				}
			}
		}
	}

	role.Menus = reqMenus

	err = rc.roleDao.UpdateRoleMenus(role)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "更新角色的权限菜单失败: "+err.Error())
		return
	}

	resp.OK(c, nil, "更新角色的权限菜单成功")

}

// 获取角色的权限接口
func (rc RoleController) GetRoleApisById(c *gin.Context) {
	// 获取path中的roleId
	roleId, _ := strconv.Atoi(c.Param("roleId"))
	if roleId <= 0 {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "角色ID不正确")
		return
	}
	// 根据path中的角色ID获取该角色信息
	roles, err := rc.roleDao.GetRolesByIds([]uint{uint(roleId)})
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, err.Error())
		return
	}
	if len(roles) == 0 {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "未获取到角色信息")
		return
	}
	// 根据角色keyword获取casbin中policy
	keyword := roles[0].Keyword
	apis, err := rc.roleDao.GetRoleApisByRoleKeyword(keyword)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, err.Error())
		return
	}
	resp.Success(c, gin.H{"apis": apis}, "获取角色的权限接口成功")
	//resp.OK(c, apis, "获取角色的权限接口成功")
}

// 更新角色的权限接口
func (rc RoleController) UpdateRoleApisById(c *gin.Context) {
	var req vo.UpdateRoleApisRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, rc.valid.Translate(err))
		return
	}

	// 获取path中的roleId
	roleId, _ := strconv.Atoi(c.Param("roleId"))
	if roleId <= 0 {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "角色ID不正确")
		return
	}
	// 根据path中的角色ID获取该角色信息
	roles, err := rc.roleDao.GetRolesByIds([]uint{uint(roleId)})
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, err.Error())
		return
	}
	if len(roles) == 0 {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "未获取到角色信息")
		return
	}

	// 当前用户角色排序最小值（最高等级角色）以及当前用户

	minSort, ctxUser, err := rc.userDao.GetCurrentUserMinRoleSort(c)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, err.Error())
		return
	}

	// (非管理员)不能更新比自己角色等级高或相等角色的权限接口
	if minSort != 1 {
		if minSort >= roles[0].Sort {
			resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "不能更新比自己角色等级高或相等角色的权限接口")
			return
		}
	}

	// 获取当前用户所拥有的权限接口
	ctxRoles := ctxUser.Roles
	ctxRolesPolicies := make([][]string, 0)
	for _, role := range ctxRoles {
		policy := rc.casbinEnforcer.GetFilteredPolicy(0, role.Keyword)
		ctxRolesPolicies = append(ctxRolesPolicies, policy...)
	}
	// 得到path中的角色ID对应角色能够设置的权限接口集合
	for _, policy := range ctxRolesPolicies {
		policy[0] = roles[0].Keyword
	}

	// 前端传来最新的ApiID集合
	apiIds := req.ApiIds
	// 根据apiID获取接口详情
	apis, err := rc.apiDao.GetApisById(apiIds)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "根据接口ID获取接口信息失败")
		return
	}
	// 生成前端想要设置的角色policies
	reqRolePolicies := make([][]string, 0)
	for _, api := range apis {
		reqRolePolicies = append(reqRolePolicies, []string{
			roles[0].Keyword, api.Path, api.Method,
		})
	}

	// (非管理员)不能把角色的权限接口设置的比当前用户所拥有的权限接口多
	if minSort != 1 {
		for _, reqPolicy := range reqRolePolicies {
			if !funk.Contains(ctxRolesPolicies, reqPolicy) {
				resp.Fail(c, ecode.RECODE_INTERNAL_ERR, fmt.Sprintf("无权设置路径为%s,请求方式为%s的接口", reqPolicy[1], reqPolicy[2]))
				return
			}
		}
	}

	// 更新角色的权限接口
	err = rc.roleDao.UpdateRoleApis(roles[0].Keyword, reqRolePolicies)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, err.Error())
		return
	}

	resp.OK(c, nil, "更新角色的权限接口成功")

}

// 批量删除角色
func (rc RoleController) BatchDeleteRoleByIds(c *gin.Context) {
	var req vo.DeleteRoleRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, rc.valid.Translate(err))
		return
	}

	// 获取当前用户最高等级角色
	minSort, _, err := rc.userDao.GetCurrentUserMinRoleSort(c)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, err.Error())
		return
	}

	// 前端传来需要删除的角色ID
	roleIds := req.RoleIds
	// 获取角色信息
	roles, err := rc.roleDao.GetRolesByIds(roleIds)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "获取角色信息失败: "+err.Error())
		return
	}
	if len(roles) == 0 {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "未获取到角色信息")
		return
	}

	// 不能删除比自己角色等级高或相等的角色
	for _, role := range roles {
		if minSort >= role.Sort {
			resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "不能删除比自己角色等级高或相等的角色")
			return
		}
	}

	// 删除角色
	err = rc.roleDao.BatchDeleteRoleByIds(roleIds)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "删除角色失败")
		return
	}

	// 删除角色成功直接清理缓存，让活跃的用户自己重新缓存最新用户信息
	rc.userDao.ClearCache() //.ClearUserInfoCache()
	resp.OK(c, nil, "删除角色成功")

}
