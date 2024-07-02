package base

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type LogOptions struct {
	Path       string        `mapstructure:"path"       json:"path"`
	Filename   string        `mapstructure:"filename"   json:"Filename"`
	MaxSize    int           `mapstructure:"max-size"    json:"maxSize"`
	MaxAge     int           `mapstructure:"max-age"     json:"maxAge"`
	MaxBackups int           `mapstructure:"max-backups" json:"maxBackups"`
	Compress   bool          `mapstructure:"compress" json:"compress"`
	Level      zapcore.Level `mapstructure:"level"    json:"level"`
	Stdout     bool          `mapstructure:"stdout"   json:"stdout"`
}

func NewLogOptions(v *viper.Viper) *LogOptions {
	var (
		err error
		o   = new(LogOptions)
	)
	if err = v.UnmarshalKey("log", o); err != nil {
		panic(err)
		return nil
	}

	return o
}

func NewLogger(o *LogOptions) *zap.Logger {
	path, file := o.Path, o.Filename
	ext := filepath.Ext(file)
	filename := strings.TrimSuffix(file, ext)

	logFileName := fmt.Sprintf("%s/%s-%04d-%02d-%02d.%s", path, filename, ext, time.Now().Year(), time.Now().Month(), time.Now().Day())
	errorLogFileName := fmt.Sprintf("%s/%s-%04d-%02d-%02d-error%s", path, filename, ext, time.Now().Year(), time.Now().Month(), time.Now().Day())

	var coreArr []zapcore.Core
	encoder := getEncoder(o.Stdout)

	// 日志级别
	highPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zap.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level < zap.ErrorLevel && level >= zap.DebugLevel
	})

	//当yml配置中的等级大于Error时，lowPriority级别日志停止记录
	if o.Level >= 2 {
		lowPriority = zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return false
		})
	}

	// info文件writeSyncer
	infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFileName,  //日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    o.MaxSize,    //文件大小限制,单位MB
		MaxAge:     o.MaxAge,     //日志文件保留天数
		MaxBackups: o.MaxBackups, //最大保留日志文件数量
		LocalTime:  false,
		Compress:   o.Compress, //是否压缩处理
	})
	// error文件writeSyncer
	errorFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   errorLogFileName, //日志文件存放目录
		MaxSize:    o.MaxSize,        //文件大小限制,单位MB
		MaxAge:     o.MaxAge,         //日志文件保留天数
		MaxBackups: o.MaxBackups,     //最大保留日志文件数量
		LocalTime:  false,
		Compress:   o.Compress, //是否压缩处理
	})

	// 第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
	infoFileCore := zapcore.NewCore(*encoder, zapcore.NewMultiWriteSyncer(infoFileWriteSyncer, zapcore.AddSync(os.Stdout)), lowPriority)

	// 第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
	errorFileCore := zapcore.NewCore(*encoder, zapcore.NewMultiWriteSyncer(errorFileWriteSyncer, zapcore.AddSync(os.Stdout)), highPriority)

	coreArr = append(coreArr, infoFileCore, errorFileCore)
	logger := zap.New(zapcore.NewTee(coreArr...), zap.AddCaller())

	zap.ReplaceGlobals(logger)
	logger.Info("初始化zap日志完成!")
	return logger
}

func InitByPathFile(pathFilename string) *zap.Logger {
	path, file := filepath.Split(pathFilename)
	o := LogOptions{
		Path:     path,
		Filename: file,
	}
	o.MaxAge = 30
	o.MaxBackups = 100
	o.MaxSize = 50
	o.Compress = false
	o.Level = -1
	o.Stdout = false
	logger := NewLogger(&o)
	return logger
}

func getEncoder(stdout bool) *zapcore.Encoder {

	// 获取编码器
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		TimeKey:       "time",
		NameKey:       "name",
		CallerKey:     "file",
		FunctionKey:   "func",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalColorLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},

		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	//if !stdout {
	//	encoder := zapcore.NewJSONEncoder(encoderConfig)
	//}
	return &encoder
}
