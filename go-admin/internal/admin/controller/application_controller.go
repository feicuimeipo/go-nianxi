package controller

import (
	"gitee.com/go-nianxi/go-common/pkg/ecode/resp"
	"gitee.com/go-nianxi/go-common/pkg/valid"
	"gitee.com/go-nianxi/go-admin/internal/admin/dao"
	"gitee.com/go-nianxi/go-admin/internal/admin/domain/dto"
	"gitee.com/go-nianxi/go-admin/internal/admin/domain/vo"
	"gitee.com/go-nianxi/go-admin/internal/admin/model"
	"gitee.com/go-nianxi/go-admin/internal/pkg/ecode"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ApplicationController struct {
	applicationDao *dao.ApplicationDao
	valid          *valid.Validator
}

func NewApplicationController(applicationDao *dao.ApplicationDao, valid *valid.Validator) *ApplicationController {
	return &ApplicationController{
		applicationDao: applicationDao,
		valid:          valid,
	}
}

// @Tags		应用
// @summary	应用列表
// @Accept		json
// @Produce	json
// @Response	200	{object}	resp.ResponseMsg
// @Param		req	body		vo.ApplicationListRequest	true	"分页查询信息"
// @Router		/app/list [get]
func (a *ApplicationController) GetApplications(c *gin.Context) {
	var req vo.ApplicationListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, a.valid.Translate(err))
		return
	}

	// 获取
	list, total, err := a.applicationDao.GetApplications(&req)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "获取应用列表失败")
		return
	}

	resp.SuccessPageList(c, list, total, "获取应用列表成功")
}

// @Tags		应用
// @summary	应用树状列表
// @Accept		json
// @Produce	json
// @Response	200	{object}	resp.ResponseMsg
// @Router		/app/tree [get]
func (a *ApplicationController) GetApplicationTree(c *gin.Context) {
	types, err := a.applicationDao.GetAllApplicationTypes()
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "获取列表失败")
		return
	}
	var treeList []dto.ApplicationTreeDto

	apps, err := a.applicationDao.GetAll()
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "获取列表失败")
		return
	}

	for _, v := range types {
		node := dto.ApplicationTreeDto{
			ID:   999900 + v.ID,
			Desc: v.Desc,
		}

		node.Children = make([]*dto.ApplicationTreeDto, 0)
		for _, av := range apps {
			if av.TypeId == v.ID {
				childNode := dto.ApplicationTreeDto{
					ID:   av.ID,
					Desc: av.Title,
				}
				node.Children = append(node.Children, &childNode)
			}
		}
		if len(node.Children) > 0 {
			treeList = append(treeList, node)
		}
	}

	resp.OK(c, treeList, "获取应用列表成功")
}

// 左则应用树立s
func (a ApplicationController) GetAppsByUserId(c *gin.Context) {
	// 获取路径中的userId
	userId, _ := strconv.Atoi(c.Param("userId"))
	if userId <= 0 {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "用户ID不正确")
		return
	}

	types, err := a.applicationDao.GetAllApplicationTypes()
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "获取列表失败")
		return
	}

	apps, err := a.applicationDao.GetAll()
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "获取用户的可访问菜单树失败: "+err.Error())
		return
	}

	var treeList []dto.ApplicationTreeDto
	for _, v := range types {
		node := dto.ApplicationTreeDto{
			ID:   999900 + v.ID,
			Desc: v.Desc,
		}

		node.Children = make([]*dto.ApplicationTreeDto, 0)
		for _, av := range apps {
			if v.ID == av.TypeId {
				childNode := dto.ApplicationTreeDto{
					ID:   av.ID,
					Desc: av.Title,
					Icon: av.Icon,
				}
				node.Children = append(node.Children, &childNode)
			}
		}
		if len(node.Children) > 0 {
			treeList = append(treeList, node)
		}
	}

	//resp.OK(c, menuTree, "获取用户的可访问菜单列表成功")
	resp.Success(c, treeList, "获取用户的可访问菜单树成功")
}

// @Tags		应用
// @summary	应用列表
// @Accept		json
// @Produce	json
// @Response	200	{object}	resp.ResponseMsg
// @Router		/app/type/list [get]
func (a *ApplicationController) GetApplicationTypes(c *gin.Context) {

	list, err := a.applicationDao.GetAllApplicationTypes()
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "获取类别列表失败")
		return
	}
	//ApplicationTypeDto
	resp.OK(c, list, "获取列表成功")
}

// @Tags		应用
// @summary	创建应用
// @Accept		json
// @Produce	json
// @Response	200	{object}	resp.ResponseMsg
// @Param		req	body		vo.CreateApplicationRequest	true	""
// @Router		/app/create [post]
func (a *ApplicationController) CreateApplication(c *gin.Context) {
	var req vo.CreateApplicationRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, a.valid.Translate(err))
		return
	}

	// 获取
	application := model.Application{
		Alias:        req.Alias,
		AppName:      req.AppName,
		Title:        req.Title,
		BaseUrl:      req.BaseUrl,
		Introduction: req.Introduction,
		Icon:         req.Icon,
	}

	// 创建接口
	err := a.applicationDao.CreateApplication(&application)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "创建失败: "+err.Error())
		return
	}

	resp.OK(c, nil, "创建接成功")
	return
}

// @Tags		应用
// @summary	更新应用
// @Accept		json
// @Produce	json
// @Response	200	{object}	resp.ResponseMsg
// @Param		req	body		vo.UpdateApplicationRequest	true	""
// @Router		/app/batch/:id [PATCH]
func (a *ApplicationController) UpdateApplicationById(c *gin.Context) {
	var req vo.UpdateApplicationRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, a.valid.Translate(err))
		return
	}

	// 获取路径中的Id
	id, _ := strconv.Atoi(c.Param("id"))
	if id <= 0 {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "接口ID不正确")
		return
	}

	// 获取当前用户
	_, err := a.applicationDao.GetApplicationById(uint(id))
	if err != nil {
		resp.Error(c, ecode.RECODE_INTERNAL_ERR, err)
		return
	}

	// 获取
	application := model.Application{
		Alias:        req.Alias,
		AppName:      req.AppName,
		Title:        req.Title,
		BaseUrl:      req.BaseUrl,
		Introduction: req.Introduction,
		Icon:         req.Icon,
	}

	err = a.applicationDao.UpdateApplicationById(uint(id), &application)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "更新失败: "+err.Error())
		return
	}

	resp.OK(c, nil, "更新成功")
}

// @Tags		应用
// @summary	批量删除
// @Accept		json
// @Produce	json
// @Response	200	{object}	resp.ResponseMsg
// @Param		req	body		vo.DeleteApplicationRequest	true	""
// @Router		/app/delete/batch [DELETE]
func (a *ApplicationController) BatchDeleteApiByIds(c *gin.Context) {
	var req vo.DeleteApplicationRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, a.valid.Translate(err))
		return
	}

	// 删除接口

	err := a.applicationDao.BatchDeleteIds(req.Ids)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "删除失败: "+err.Error())
		return
	}

	resp.OK(c, nil, "删除成功")
}
