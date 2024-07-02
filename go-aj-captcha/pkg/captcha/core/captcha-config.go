package core

import (
	"errors"
	"gitee.com/go-nianxi/go-aj-captcha/pkg/captcha/cache"
	"image/color"
	"strings"
)

type CaptchaConfig struct {
	Watermark   *WatermarkConfig
	ClickWord   *ClickWordConfig
	BlockPuzzle *BlockPuzzleConfig
	// 验证码使用的缓存类型
	CacheType      string
	CacheExpireSec int
	// 项目的绝对路径: 图片、字体等
	ResourcePath string
}

// WatermarkConfig 水印设置
type WatermarkConfig struct {
	FontSize int        `yaml:"fontSize"`
	Color    color.RGBA `yaml:"color"`
	Text     string     `yaml:"text"`
}

type BlockPuzzleConfig struct {
	// 校验时 容错偏移量
	Offset int `yaml:"offset"`
}

type ClickWordConfig struct {
	FontSize int `yaml:"fontSize"`
	FontNum  int `yaml:"fontNum"`
}

func NewCaptchaConfig() *CaptchaConfig {
	return &CaptchaConfig{
		//可以为redis类型缓存RedisCacheKey，也可以为内存MemCacheKey
		CacheType: cache.MemCacheKey,
		Watermark: &WatermarkConfig{
			FontSize: 12,
			Color:    color.RGBA{R: 255, G: 255, B: 255, A: 255},
			Text:     "我的水印",
		},
		ClickWord: &ClickWordConfig{
			FontSize: 25,
			FontNum:  4,
		},
		BlockPuzzle:    &BlockPuzzleConfig{Offset: 10},
		CacheExpireSec: 2 * 60, // 缓存有效时间
		ResourcePath:   "./",
	}
}

// BuildConfig 生成config配置
func InitCaptchaConfig(
	cacheType,
	resourcePath string,
	waterConfig *WatermarkConfig,
	clickWordConfig *ClickWordConfig,
	puzzleConfig *BlockPuzzleConfig,
	cacheExpireSec int,
) *CaptchaConfig {

	if len(resourcePath) == 0 {
		resourcePath = DefaultResourceRoot
	}

	if len(cacheType) == 0 {
		cacheType = cache.MemCacheKey
	} else if strings.Compare(cacheType, cache.MemCacheKey) != 0 && strings.Compare(cacheType, cache.RedisCacheKey) != 0 {
		panic(errors.New("cache type not support"))
		return nil
	}
	if cacheExpireSec == 0 {
		cacheExpireSec = 2 * 60
	}
	if nil == waterConfig {
		waterConfig = &WatermarkConfig{
			FontSize: 12,
			Color:    color.RGBA{R: 255, G: 255, B: 255, A: 255},
			Text:     DefaultText,
		}
	}
	if nil == clickWordConfig {
		clickWordConfig = &ClickWordConfig{
			FontSize: 25,
			FontNum:  4,
		}
	}
	if nil == puzzleConfig {
		puzzleConfig = &BlockPuzzleConfig{Offset: 10}
	}

	return &CaptchaConfig{
		//可以为redis类型缓存RedisCacheKey，也可以为内存MemCacheKey
		CacheType:   cacheType,
		Watermark:   waterConfig,
		ClickWord:   clickWordConfig,
		BlockPuzzle: puzzleConfig,
		// 缓存有效时间
		CacheExpireSec: cacheExpireSec,
		ResourcePath:   resourcePath,
	}
}
