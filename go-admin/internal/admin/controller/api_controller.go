package controller

import (
	"gitee.com/go-nianxi/go-common/pkg/ecode/resp"
	"gitee.com/go-nianxi/go-common/pkg/valid"
	"gitee.com/go-nianxi/go-admin/internal/admin/dao"
	"gitee.com/go-nianxi/go-admin/internal/admin/domain/vo"
	"gitee.com/go-nianxi/go-admin/internal/admin/model"
	"gitee.com/go-nianxi/go-admin/internal/pkg/ecode"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ApiController struct {
	apiDao  *dao.ApiDao
	userDao *dao.UserDao
	appDao  *dao.ApplicationDao
	valid   *valid.Validator
}

func NewApiController(apiDao *dao.ApiDao, userDao *dao.UserDao, appDao *dao.ApplicationDao, valid *valid.Validator) *ApiController {
	return &ApiController{
		apiDao:  apiDao,
		userDao: userDao,
		valid:   valid,
		appDao:  appDao,
	}
}

// @Tags		接口
// @summary	获取接口列表
// @Accept		json
// @Produce	json
// @Response	200	{object}	resp.ResponseMsg
// @Param		req	body		vo.ApiListRequest	true	"登录信息"
// @Router		/api/list [get]
//
// @Security	ApiKeyAuth
func (ac *ApiController) GetApis(c *gin.Context) {
	var req vo.ApiListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, ac.valid.Translate(err))
		return
	}

	if req.ApplicationId == 0 {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "1.应用编号不可以为空！")
		return
	}

	// 获取
	apis, total, err := ac.apiDao.GetApis(&req)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "获取接口列表失败")
		return
	}

	//resp.OK(c, resp.PageData{List: apis, Total: total}, "获取接口树成功")
	resp.Success(c, gin.H{
		"apis": apis, "total": total,
	}, "获取接口列表成功")
}

// 获取接口树(按接口Category字段分类)
func (ac *ApiController) GetApiTree(c *gin.Context) {
	// 获取路径中的apiId
	appId, _ := strconv.Atoi(c.Param("appId"))
	if appId <= 0 {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "应用ID不正确")
		return
	}

	tree, err := ac.apiDao.GetApiTree(uint(appId))
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "获取接口树失败")
		return
	}

	//resp.OK(c, tree, "获取接口树成功")
	resp.Success(c, gin.H{
		"apiTree": tree,
	}, "获取接口树成功")
}

// 创建接口
func (ac *ApiController) CreateApi(c *gin.Context) {
	var req vo.CreateApiRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, ac.valid.Translate(err))
		return
	}

	// 获取当前用户
	ctxUser, err := ac.userDao.GetCurrentUser(c)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "获取当前用户信息失败")
		return
	}

	app, err := ac.appDao.GetApplicationById(req.ApplicationId)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, err.Error())
		return
	}

	api := model.Api{
		Method:        req.Method,
		Path:          req.Path,
		Category:      req.Category,
		Desc:          req.Desc,
		ApplicationId: app.ID,
		Creator:       ctxUser.Username,
	}

	// 创建接口
	err = ac.apiDao.CreateApi(&api)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "创建接口失败: "+err.Error())
		return
	}

	resp.OK(c, nil, "创建接口成功")
	return
}

// 更新接口
func (ac *ApiController) UpdateApiById(c *gin.Context) {
	var req vo.UpdateApiRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, ac.valid.Translate(err))
		return
	}

	// 获取路径中的apiId
	apiId, _ := strconv.Atoi(c.Param("apiId"))
	if apiId <= 0 {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "接口ID不正确")
		return
	}

	// 获取当前用户
	ctxUser, err := ac.userDao.GetCurrentUser(c)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "获取当前用户信息失败")
		return
	}

	api := model.Api{
		Method:   req.Method,
		Path:     req.Path,
		Category: req.Category,
		Desc:     req.Desc,
		Creator:  ctxUser.Username,
	}

	err = ac.apiDao.UpdateApiById(uint(apiId), &api)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "更新接口失败: "+err.Error())
		return
	}

	resp.OK(c, nil, "更新接口成功")
}

// 批量删除接口
func (ac *ApiController) BatchDeleteApiByIds(c *gin.Context) {
	var req vo.DeleteApiRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, ac.valid.Translate(err))
		return
	}

	// 删除接口

	err := ac.apiDao.BatchDeleteApiByIds(req.ApiIds)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "删除接口失败: "+err.Error())
		return
	}

	resp.OK(c, nil, "删除接口成功")
}
