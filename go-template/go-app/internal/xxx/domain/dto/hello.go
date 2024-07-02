package dto

import "gitee.com/go-nianxi/go-template/internal/xxx/model"

type HelloDTO struct {
	Msg string `gorm:"type:varchar(20);comment:'hello'" json:"hello"`
}

func ToHelloDTO(model *model.Hello) *HelloDTO {
	return &HelloDTO{
		Msg: model.Msg,
	}

}
