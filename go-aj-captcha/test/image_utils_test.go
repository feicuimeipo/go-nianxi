package test

import (
	"gitee.com/go-nianxi/go-aj-captcha/pkg/captcha/core"
	"gitee.com/go-nianxi/go-aj-captcha/pkg/captcha/util"
	"testing"
)

func TestImage_GetBackgroundImage(t *testing.T) {
	background := util.NewImageByBackgroundImage(core.DefaultFont)
	template := util.NewImageByTemplateImage(core.DefaultFont)

	if background == nil {
		t.Fatal("背景图片获取失败")
	}
	if template == nil {
		t.Fatal("模板图片获取失败")
	}
}
