package controller

import (
	"gitee.com/go-nianxi/go-common/pkg/ecode"
	"gitee.com/go-nianxi/go-common/pkg/ecode/resp"
	"gitee.com/go-nianxi/go-common/pkg/valid"
	"gitee.com/go-nianxi/go-template/internal/admin/domain/vo"
	"gitee.com/go-nianxi/go-template/internal/admin/repository"
	"github.com/gin-gonic/gin"
)

type HelloController struct {
	helloRepository *repository.HelloRepository
	valid           *valid.Validator
}

func NewHelloController(helloDao *repository.HelloRepository, valid *valid.Validator) *HelloController {
	return &HelloController{
		helloRepository: helloDao,
		valid:           valid,
	}
}

func (h *HelloController) GetHellos(c *gin.Context) {
	var req vo.HelloListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, h.valid.Translate(err))
		return
	}

	list, total, err := h.helloRepository.GetHellos(&req)
	if err != nil {
		resp.Fail(c, ecode.RECODE_INTERNAL_ERR, "获取操作日志列表失败: "+err.Error())
		return
	}

	resp.SuccessPageList(c, list, total)
}
