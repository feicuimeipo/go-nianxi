package test

import (
	"fmt"
	"gitee.com/go-nianxi/go-aj-captcha/pkg/captcha/util"
	"testing"
)

func TestRandomInt(t *testing.T) {
	fmt.Println(util.RandomInt(198, 236))
}
