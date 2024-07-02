package controller

import (
	"gitee.com/go-nianxi/go-common/pkg/ecode/resp"
	"gitee.com/go-nianxi/go-common/pkg/valid"
	"gitee.com/go-nianxi/go-admin/internal/admin/dao"
	"gitee.com/go-nianxi/go-admin/internal/admin/domain/vo"
	"gitee.com/go-nianxi/go-admin/internal/pkg/ecode"
	"github.com/gin-gonic/gin"
)

type OperationLogController struct {
	operationDao *dao.OperationLogDao
	valid        *valid.Validator
}

func NewOperationLogController(operationDao *dao.OperationLogDao, valid *valid.Validator) *OperationLogController {
	return &OperationLogController{
		operationDao: operationDao,
		valid:        valid,
	}
}

// 获取操作日志列表
func (oc OperationLogController) GetOperationLogs(c *gin.Context) {
	var req vo.OperationLogListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, oc.valid.Translate(err))
		return
	}
	// 获取
	logs, total, err := oc.operationDao.GetOperationLogs(&req)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "获取操作日志列表失败: "+err.Error())
		return
	}
	resp.Success(c, gin.H{"logs": logs, "total": total}, "获取操作日志列表成功")
	//resp.OK(c, resp.PageData{List: logs, Total: total}, "获取操作日志列表成功")
}

// 批量删除操作日志
func (oc OperationLogController) BatchDeleteOperationLogByIds(c *gin.Context) {
	var req vo.DeleteOperationLogRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, oc.valid.Translate(err))
		return
	}

	// 删除接口
	err := oc.operationDao.BatchDeleteOperationLogByIds(req.OperationLogIds)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "删除日志失败: "+err.Error())
		return
	}
	resp.OK(c, nil, "删除日志成功")
	//response.Success(c, nil, "删除日志成功")
}
