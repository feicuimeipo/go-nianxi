package test

import (
	"fmt"
	"gitee.com/go-nianxi/go-aj-captcha/pkg/captcha/cache"
	"gitee.com/go-nianxi/go-aj-captcha/pkg/captcha/util"
	"image/color"
	"testing"
)

func TestBlockPuzzleCaptchaService_Get(t *testing.T) {
	//
	//vo := &vo2.CaptchaVO{}
	//b := &core.BlockPuzzleCaptchaService{}
	//res := b.Get(*vo)
	//
	//fmt.Println(res)
}

func TestImage(t *testing.T) {

	backgroundImage := util.NewImageUtil("/mnt/f/workspace/core/resources/defaultImages/jigsaw/original/1.png", "/mnt/f/workspace/core/resources/fonts/WenQuanZhengHei.ttf")
	// 为背景图片设置水印
	backgroundImage.SetText("牛逼AAA", 14, color.RGBA{R: 120, G: 120, B: 255, A: 255})
	backgroundImage.DecodeImageToFile()
}

func TestIntCovert(t *testing.T) {

	cache := cache.NewMemCacheService(10)

	cache.Set("test1", "tes111", 0)

	val := cache.Get("test1")

	fmt.Println(val)

}
