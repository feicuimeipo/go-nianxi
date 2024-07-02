package controller

import (
	"gitee.com/go-nianxi/go-common/pkg/base"
	"gitee.com/go-nianxi/go-common/pkg/ecode/resp"
	"gitee.com/go-nianxi/go-common/pkg/transports/http"
	"gitee.com/go-nianxi/go-common/pkg/valid"
	"gitee.com/go-nianxi/go-admin/internal/admin/config"
	"gitee.com/go-nianxi/go-admin/internal/admin/dao"
	"gitee.com/go-nianxi/go-admin/internal/admin/domain/dto"
	"gitee.com/go-nianxi/go-admin/internal/admin/domain/vo"
	"gitee.com/go-nianxi/go-admin/internal/admin/model"
	"gitee.com/go-nianxi/go-admin/internal/admin/util"
	"gitee.com/go-nianxi/go-admin/internal/pkg/ecode"
	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
	"strconv"
)

type UserController struct {
	userDao *dao.UserDao
	roleDao *dao.RoleDao
	valid   *valid.Validator
}

// 构造函数
func NewUserController(userDao *dao.UserDao, roleDao *dao.RoleDao, valid *valid.Validator) *UserController {
	return &UserController{
		userDao: userDao,
		valid:   valid,
		roleDao: roleDao,
	}
}

// 获取当前登录用户信息
func (uc *UserController) GetUserInfo(c *gin.Context) {
	user, err := uc.userDao.GetCurrentUser(c)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "获取当前用户信息失败: "+err.Error())
		return
	}
	userInfoDto := dto.ToUserInfoDto(user)

	//resp.OK(c, userInfoDto, "获取当前用户信息成功！")
	resp.Success(c, gin.H{
		"userInfo": userInfoDto,
	}, "获取当前用户信息成功")
}

// 获取用户列表
func (uc *UserController) GetUsers(c *gin.Context) {
	var req vo.UserListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, uc.valid.Translate(err))
		return
	}

	// 获取
	users, total, err := uc.userDao.GetUsers(&req)
	if err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, "获取用户列表失败: "+err.Error())
		return
	}

	resp.Success(c, gin.H{"users": dto.ToUsersDto(users), "total": total}, "获取用户列表成功")
	//resp.Success(c, resp.PageData{List: users, Total: total})
}

// 更新用户登录密码
func (uc *UserController) ChangePwd(c *gin.Context) {
	var req vo.ChangePwdRequest

	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, uc.valid.Translate(err))
		return
	}

	// 前端传来的密码是rsa加密的,先解密
	// 密码通过RSA解密
	if config.Conf.O.Http.Mode == base.PROD_MODE || len(req.OldPassword) > 20 {
		decodeOldPassword, err := http.RSADecrypt([]byte(req.OldPassword), config.Conf.O.Http.Ssl.KeyFileBytes)
		if err != nil {
			resp.Fail(c, ecode.RECODE_INTERNAL_ERR, err.Error())
			return
		}
		decodeNewPassword, err := http.RSADecrypt([]byte(req.NewPassword), config.Conf.O.Http.Ssl.KeyFileBytes)
		if err != nil {
			resp.Fail(c, ecode.RECODE_INTERNAL_ERR, err.Error())
			return
		}
		req.OldPassword = string(decodeOldPassword)
		req.NewPassword = string(decodeNewPassword)
	}

	// 获取当前用户
	user, err := uc.userDao.GetCurrentUser(c)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, err.Error())
		return
	}
	// 获取用户的真实正确密码
	correctPasswd := user.Password
	// 判断前端请求的密码是否等于真实密码
	err = util.ComparePasswd(correctPasswd, req.OldPassword)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "原密码有误")
		return
	}
	// 更新密码
	err = uc.userDao.ChangePwd(user.Username, util.GenPasswd(req.NewPassword))
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "更新密码失败: "+err.Error())
		return
	}
	resp.OK(c, ecode.RECODE_INTERNAL_ERR, "更新密码成功")
}

// 创建用户
func (uc *UserController) CreateUser(c *gin.Context) {
	var req vo.CreateUserRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, uc.valid.Translate(err))
		return
	}

	// 密码通过RSA解密
	// 密码不为空就解密
	if req.Password != "" {
		decodeData, err := http.RSADecrypt([]byte(req.Password), config.Conf.O.Http.Ssl.KeyFileBytes)
		if err != nil {
			resp.Fail(c, ecode.RECODE_INTERNAL_ERR, err.Error())
			return
		}
		req.Password = string(decodeData)
		if len(req.Password) < 6 {
			resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "密码长度至少为6位")
			return
		}
	}

	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	currentRoleSortMin, ctxUser, err := uc.userDao.GetCurrentUserMinRoleSort(c)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, err.Error())
		return
	}

	// 获取前端传来的用户角色id
	reqRoleIds := req.RoleIds
	// 根据角色id获取角色

	roles, err := uc.roleDao.GetRolesByIds(reqRoleIds)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "根据角色ID获取角色信息失败: "+err.Error())
		return
	}
	if len(roles) == 0 {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "未获取到角色信息")
		return
	}
	var reqRoleSorts []int
	for _, role := range roles {
		reqRoleSorts = append(reqRoleSorts, int(role.Sort))
	}
	// 前端传来用户角色排序最小值（最高等级角色）
	reqRoleSortMin := uint(funk.MinInt(reqRoleSorts))

	// 当前用户的角色排序最小值 需要小于 前端传来的角色排序最小值（用户不能创建比自己等级高的或者相同等级的用户）
	if currentRoleSortMin >= reqRoleSortMin {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "用户不能创建比自己等级高的或者相同等级的用户")
		return
	}

	// 密码为空就默认123456
	if req.Password == "" {
		req.Password = "123456"
	}
	user := model.User{
		Username:     req.Username,
		Password:     util.GenPasswd(req.Password),
		Mobile:       req.Mobile,
		Avatar:       req.Avatar,
		Nickname:     &req.Nickname,
		Introduction: &req.Introduction,
		Status:       req.Status,
		Creator:      ctxUser.Username,
		Roles:        roles,
	}

	err = uc.userDao.CreateUser(&user)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "创建用户失败: "+err.Error())
		return
	}
	resp.OK(c, ecode.RECODE_INTERNAL_ERR, "创建用户成功")

}

// 更新用户
func (uc *UserController) UpdateUserById(c *gin.Context) {
	var req vo.CreateUserRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, uc.valid.Translate(err))
		return
	}

	//获取path中的userId
	userId, _ := strconv.Atoi(c.Param("userId"))
	if userId <= 0 {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "用户ID不正确")
		return
	}

	// 根据path中的userId获取用户信息
	oldUser, err := uc.userDao.GetUserById(uint(userId))
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "获取需要更新的用户信息失败: "+err.Error())
		return
	}

	// 获取当前用户
	ctxUser, err := uc.userDao.GetCurrentUser(c)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, err.Error())
		return
	}
	// 获取当前用户的所有角色
	currentRoles := ctxUser.Roles
	// 获取当前用户角色的排序，和前端传来的角色排序做比较
	var currentRoleSorts []int
	// 当前用户角色ID集合
	var currentRoleIds []uint
	for _, role := range currentRoles {
		currentRoleSorts = append(currentRoleSorts, int(role.Sort))
		currentRoleIds = append(currentRoleIds, role.ID)
	}
	// 当前用户角色排序最小值（最高等级角色）
	currentRoleSortMin := funk.MinInt(currentRoleSorts)

	// 获取前端传来的用户角色id
	reqRoleIds := req.RoleIds
	// 根据角色id获取角色

	roles, err := uc.roleDao.GetRolesByIds(reqRoleIds)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "根据角色ID获取角色信息失败: "+err.Error())
		return
	}
	if len(roles) == 0 {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "未获取到角色信息")
		return
	}
	var reqRoleSorts []int
	for _, role := range roles {
		reqRoleSorts = append(reqRoleSorts, int(role.Sort))
	}
	// 前端传来用户角色排序最小值（最高等级角色）
	reqRoleSortMin := funk.MinInt(reqRoleSorts)

	user := model.User{
		Model:        oldUser.Model,
		Username:     req.Username,
		Password:     oldUser.Password,
		Mobile:       req.Mobile,
		Avatar:       req.Avatar,
		Nickname:     &req.Nickname,
		Introduction: &req.Introduction,
		Status:       req.Status,
		Creator:      ctxUser.Username,
		Roles:        roles,
	}
	// 判断是更新自己还是更新别人
	if userId == int(ctxUser.ID) {
		// 如果是更新自己
		// 不能禁用自己
		if req.Status == 2 {
			resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "不能禁用自己")
			return
		}
		// 不能更改自己的角色
		reqDiff, currentDiff := funk.Difference(req.RoleIds, currentRoleIds)
		if len(reqDiff.([]uint)) > 0 || len(currentDiff.([]uint)) > 0 {
			resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "不能更改自己的角色")
			return
		}

		// 不能更新自己的密码，只能在个人中心更新
		if req.Password != "" {
			resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "请到个人中心更新自身密码")
			return
		}

		// 密码赋值
		user.Password = ctxUser.Password

	} else {
		// 如果是更新别人
		// 用户不能更新比自己角色等级高的或者相同等级的用户
		// 根据path中的userIdID获取用户角色排序最小值
		minRoleSorts, err := uc.userDao.GetUserMinRoleSortsByIds([]uint{uint(userId)})
		if err != nil || len(minRoleSorts) == 0 {
			resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "根据用户ID获取用户角色排序最小值失败")
			return
		}
		if currentRoleSortMin >= minRoleSorts[0] {
			resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "用户不能更新比自己角色等级高的或者相同等级的用户")
			return
		}

		// 用户不能把别的用户角色等级更新得比自己高或相等
		if currentRoleSortMin >= reqRoleSortMin {
			resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "用户不能把别的用户角色等级更新得比自己高或相等")
			return
		}

		// 密码赋值
		if req.Password != "" {
			// 密码通过RSA解密
			decodeData, err := http.RSADecrypt([]byte(req.Password), config.Conf.O.Http.Ssl.KeyFileBytes)
			if err != nil {
				resp.Fail(c, ecode.RECODE_INTERNAL_ERR, err.Error())
				return
			}
			req.Password = string(decodeData)
			user.Password = util.GenPasswd(req.Password)
		}

	}

	// 更新用户
	err = uc.userDao.UpdateUser(&user)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "更新用户失败: "+err.Error())
		return
	}
	resp.OK(c, ecode.RECODE_INTERNAL_ERR, "更新用户成功")

}

// 批量删除用户
func (uc *UserController) BatchDeleteUserByIds(c *gin.Context) {
	var req vo.DeleteUserRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, uc.valid.Translate(err))
		return
	}

	// 前端传来的用户ID
	reqUserIds := req.UserIds
	// 根据用户ID获取用户角色排序最小值
	roleMinSortList, err := uc.userDao.GetUserMinRoleSortsByIds(reqUserIds)
	if err != nil || len(roleMinSortList) == 0 {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "根据用户ID获取用户角色排序最小值失败")
		return
	}

	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	minSort, ctxUser, err := uc.userDao.GetCurrentUserMinRoleSort(c)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, err.Error())
		return
	}
	currentRoleSortMin := int(minSort)

	// 不能删除自己
	if funk.Contains(reqUserIds, ctxUser.ID) {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "用户不能删除自己")
		return
	}

	// 不能删除比自己角色排序低(等级高)的用户
	for _, sort := range roleMinSortList {
		if currentRoleSortMin >= sort {
			resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "用户不能删除比自己角色等级高的用户")
			return
		}
	}

	err = uc.userDao.BatchDeleteUserByIds(reqUserIds)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "删除用户失败: "+err.Error())
		return
	}

	resp.OK(c, ecode.RECODE_INTERNAL_ERR, "删除用户成功")

}
