package controller

import (
	"gitee.com/go-nianxi/go-common/pkg/ecode/resp"
	"gitee.com/go-nianxi/go-common/pkg/valid"
	"gitee.com/go-nianxi/go-admin/internal/admin/dao"
	"gitee.com/go-nianxi/go-admin/internal/admin/domain/vo"
	"gitee.com/go-nianxi/go-admin/internal/admin/model"
	"gitee.com/go-nianxi/go-admin/internal/pkg/ecode"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

type MenuController struct {
	menuDao *dao.MenuDao
	userDao *dao.UserDao
	appDao  *dao.ApplicationDao
	valid   *valid.Validator
	logger  *zap.Logger
}

func NewMenuController(dao *dao.Dao, valid *valid.Validator, logger *zap.Logger) *MenuController {

	return &MenuController{
		menuDao: dao.MenuDao,
		userDao: dao.UserDao,
		valid:   valid,
		appDao:  dao.ApplicationDao,
		logger:  logger,
	}
}

// 获取菜单列表
func (mc MenuController) GetMenus(c *gin.Context) {
	// 获取路径中的apiId
	appId, _ := strconv.Atoi(c.Param("appId"))
	if appId <= 0 {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "应用编号ID不正确")
		return
	}

	menus, err := mc.menuDao.GetMenus(uint(appId))
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "获取菜单列表失败: "+err.Error())
		return
	}
	resp.Success(c, gin.H{"menus": menus}, "获取菜单列表成功")
	//resp.OK(c, menus, "获取菜单列表成功")
}

// 获取菜单树
func (mc MenuController) GetMenuTree(c *gin.Context) {
	// 获取路径中的apiId
	appId, _ := strconv.Atoi(c.Param("appId"))
	if appId <= 0 {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "应用ID不正确")
		return
	}

	menuTree, err := mc.menuDao.GetMenuTree(uint(appId))
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "获取菜单树失败: "+err.Error())
		return
	}
	resp.Success(c, gin.H{"menuTree": menuTree}, "获取菜单树成功")
	//resp.OK(c, menuTree, "获取菜单列表成功")
}

// 创建菜单
func (mc MenuController) CreateMenu(c *gin.Context) {
	var req vo.CreateMenuRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, mc.valid.Translate(err))
		return
	}

	// 获取当前用户
	//ur := repository.NewUserRepository()
	ctxUser, err := mc.userDao.GetCurrentUser(c)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "获取当前用户信息失败")
		return
	}

	app, err := mc.appDao.GetApplicationById(req.ApplicationId)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, err.Error())
		return
	}

	menu := model.Menu{
		Name:          req.Name,
		Title:         req.Title,
		Icon:          &req.Icon,
		Path:          req.Path,
		Redirect:      &req.Redirect,
		Component:     req.Component,
		Sort:          req.Sort,
		Status:        req.Status,
		Hidden:        req.Hidden,
		NoCache:       req.NoCache,
		AlwaysShow:    req.AlwaysShow,
		Breadcrumb:    req.Breadcrumb,
		ActiveMenu:    &req.ActiveMenu,
		ParentId:      &req.ParentId,
		Creator:       ctxUser.Username,
		ApplicationId: app.ID,
	}

	err = mc.menuDao.CreateMenu(&menu)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "创建菜单失败: "+err.Error())
		return
	}
	resp.OK(c, nil, "创建菜单成功")
}

// 更新菜单
func (mc MenuController) UpdateMenuById(c *gin.Context) {
	var req vo.UpdateMenuRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, mc.valid.Translate(err))
		return
	}

	// 获取路径中的menuId
	menuId, _ := strconv.Atoi(c.Param("menuId"))
	if menuId <= 0 {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "菜单ID不正确")
		return
	}

	// 获取当前用户
	//ur := repository.NewUserRepository()
	ctxUser, err := mc.userDao.GetCurrentUser(c)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "获取当前用户信息失败")
		return
	}

	menu := model.Menu{
		Name:          req.Name,
		Title:         req.Title,
		Icon:          &req.Icon,
		Path:          req.Path,
		Redirect:      &req.Redirect,
		Component:     req.Component,
		Sort:          req.Sort,
		Status:        req.Status,
		Hidden:        req.Hidden,
		NoCache:       req.NoCache,
		AlwaysShow:    req.AlwaysShow,
		Breadcrumb:    req.Breadcrumb,
		ActiveMenu:    &req.ActiveMenu,
		ParentId:      &req.ParentId,
		ApplicationId: req.ApplicationId,
		Creator:       ctxUser.Username,
	}

	err = mc.menuDao.UpdateMenuById(uint(menuId), &menu)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "更新菜单失败: "+err.Error())
		return
	}

	resp.OK(c, nil, "更新菜单成功")

}

// 批量删除菜单
func (mc MenuController) BatchDeleteMenuByIds(c *gin.Context) {
	var req vo.DeleteMenuRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, mc.valid.Translate(err))
		return
	}
	err := mc.menuDao.BatchDeleteMenuByIds(req.MenuIds)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "删除菜单失败: "+err.Error())
		return
	}

	resp.OK(c, nil, "删除菜单成功")
}

// 目前用途未知
// 根据用户ID获取用户的可访问菜单列表
func (mc MenuController) GetUserMenusByUserId(c *gin.Context) {
	// 获取路径中的userId
	userId, _ := strconv.Atoi(c.Param("userId"))
	if userId <= 0 {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "用户ID不正确")
		return
	}

	menus, err := mc.menuDao.GetUserMenusByUserId(uint(userId), 0)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "获取用户的可访问菜单列表失败: "+err.Error())
		return
	}
	resp.Success(c, gin.H{"menus": menus}, "获取用户的可访问菜单列表成功")
	//resp.OK(c, menus, "获取用户的可访问菜单列表成功")
}

// 左则菜单树
// 根据用户ID及菜单Id获取用户的可访问菜单树
func (mc MenuController) GetUserMenuTreeByUserId(c *gin.Context) {
	// 获取路径中的userId
	userId, _ := strconv.Atoi(c.Param("userId"))
	if userId <= 0 {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "用户ID不正确")
		return
	}

	// 获取路径中的userId
	appId, _ := strconv.Atoi(c.Param("appId"))
	if appId <= 0 {
		mc.logger.Error("应用id没有被正确传递，给默认值1")
		appId = 1
	}

	menuTree, err := mc.menuDao.GetUserMenuTreeByUserId(uint(userId), uint(appId))
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "获取用户的可访问菜单树失败: "+err.Error())
		return
	}
	//resp.OK(c, menuTree, "获取用户的可访问菜单列表成功")
	resp.Success(c, gin.H{"menuTree": menuTree}, "获取用户的可访问菜单树成功")
}
