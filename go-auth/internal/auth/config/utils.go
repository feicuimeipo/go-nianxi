package config

import (
	"os"
	"path"
	"runtime"
)

// 获取当前执行文件绝对路径（go run）
func getCurrentAbPathByCaller1() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

func getCurrentAbPathByCaller() string {
	filename, _ := os.Getwd()
	return filename
}
