package core

import (
	"image/color"
)

const (
	// CodeKeyPrefix 缓存key前缀
	CodeKeyPrefix = "RUNNING:CAPTCHA:%s"

	// BlockPuzzleCaptcha 滑动验证码服务标识
	BlockPuzzleCaptcha = "blockPuzzle"

	// ClickWordCaptcha 点击验证码服务标识
	ClickWordCaptcha = "clickWord"

	// DefaultFont 字体文件地址
	DefaultFontFile = "/resources/fonts/WenQuanZhengHei.ttf"

	// DefaultText 默认水印显示文字
	DefaultText = "我的水印"

	// DefaultResourceRoot 默认根目录
	DefaultResourceRoot = ""
)

const (
	// DefaultTemplateImageDirectory 滑动模板图文件目录地址
	DefaultTemplateImageDirectory = "/resources/defaultImages/jigsaw/slidingBlock"
	// DefaultBackgroundImageDirectory 背景图片目录地址
	DefaultBackgroundImageDirectory = "/resources/defaultImages/jigsaw/original"
	// DefaultClickBackgroundImageDirectory 点击背景图默认地址
	DefaultClickBackgroundImageDirectory = "/resources/defaultImages/pic-click"
)

// **********************默认配置***************************************************
// var config = captchaConfig.NewConfig()
// *********************自定义配置**************************************************
// 水印配置（参数可从业务系统自定义）
var watermarkConfig = &WatermarkConfig{
	FontSize: 12,
	Color:    color.RGBA{R: 255, G: 255, B: 255, A: 255},
	Text:     "我的水印",
}

// 点击文字配置（参数可从业务系统自定义）
var clickWordConfig = &ClickWordConfig{
	FontSize: 25,
	FontNum:  4,
}

// 滑动模块配置（参数可从业务系统自定义）
var blockPuzzleConfig = &BlockPuzzleConfig{Offset: 10}
