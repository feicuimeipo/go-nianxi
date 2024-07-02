package base

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

const (
	UrlPathPrefix = "/"
	PROD_MODE     = "release"
	SystemUser    = 1
)

var Conf = new(BaseConf)

type BaseConf struct {
	O      *AppOptions
	Viper  *viper.Viper
	Logger *zap.Logger
}

type AppOptions struct {
	Mode     string `mapstructure:"mode"`
	Name     string `mapstructure:"name"`
	InitData bool   `mapstructure:"init-data"`
	
	IsAdminUser       bool        `mapstructure:"-"`
	IsProd            bool        `mapstructure:"-"`
	DefaultTenantId   uint        `mapstructure:"-"                 json:"-"`
	DefaultTenantCode string      `mapstructure:"-"                 json:"-"`
	Log               *LogOptions `mapstructure:"-"                 json:"-"`
}

func NewAppBaseOptions(v *viper.Viper) *AppOptions {
	var err error
	o := new(AppOptions)
	key := "app"
	if err = v.UnmarshalKey(key, o); err != nil {
		panic(fmt.Errorf("初始化 %s 配置失败:%s \n", key, err))
		return nil
	}

	o.IsProd = o.Mode == PROD_MODE
	o.DefaultTenantId = 1
	o.DefaultTenantCode = "default"

	o.Log = NewLogOptions(v)
	return o
}

func InitViper(configFile string) *viper.Viper {
	var v = viper.New()

	path, fileName := filepath.Split(configFile)
	//ext := filepath.Ext(fileName)
	if !filepath.IsAbs(path) {
		workDir, err := os.Getwd()
		if err != nil {
			panic(fmt.Errorf("读取应用目录失败:%s \n", err))

		}
		path = filepath.Join(workDir, "configs", path)
		v.SetConfigFile(path + "\\" + fileName)
	} else {
		v.SetConfigFile(configFile)
	}

	v.SetConfigType("yml")

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("\n 配置文件路径 -> %s\\%s \n 加载失败 -> %v\n", path, fileName, v))
		return nil

	}

	//viper.WatchConfig()
	fmt.Printf("\n 配置文件路径 -> %s\\%s \n 加载成功 -> %s\n", path, fileName, v.ConfigFileUsed())

	//以下为初始化数据
	return v
}

// NewAppBase 有错直接抛出异常
func NewAppBase(viper *viper.Viper) *BaseConf {
	Conf.Viper = viper
	Conf.O = NewAppBaseOptions(viper)
	Conf.Logger = NewLogger(Conf.O.Log)
	return Conf
}
