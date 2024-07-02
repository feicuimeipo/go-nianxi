package router

import (
	"gitee.com/go-nianxi/go-common/pkg/ecode"
	"gitee.com/go-nianxi/go-common/pkg/ecode/resp"
	"gitee.com/go-nianxi/go-common/pkg/valid"
	dto2 "gitee.com/go-nianxi/go-template/internal/xxx/domain/dto"
	"gitee.com/go-nianxi/go-template/internal/xxx/domain/vo"
	"gitee.com/go-nianxi/go-template/internal/xxx/service"
	"github.com/gin-gonic/gin"
)

type HelloRoutes struct {
	helloService *service.HelloService
	valid        *valid.Validator
}

func InitHelloRoutes(r *gin.RouterGroup, service *service.Service, valid *valid.Validator) {
	h := new(HelloRoutes)
	h.helloService = service.HelloService
	h.valid = valid

	r.GET("/hello", h.GetHelloById)

}

func (h *HelloRoutes) GetHelloById(c *gin.Context) {
	var req vo.HelloRequest
	if err := c.ShouldBind(&req); err != nil {
		resp.Fail(c, ecode.RECODE_REQ_PARAM_ERR, h.valid.Translate(err))
		return
	}

	hello, err := h.helloService.GetHelloById(&req)
	if err != nil {
		resp.Error(c, ecode.RECODE_INTERNAL_ERR, err)
	}

	dto := dto2.ToHelloDTO(hello)
	resp.Success(c, dto)

}
