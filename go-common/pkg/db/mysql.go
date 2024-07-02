package db

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
)

// 全局mysql数据库变量
//var DB *gorm.DB

type MysqlOptions struct {
	Username    string `mapstructure:"username" json:"username"`
	Password    string `mapstructure:"password" json:"password"`
	Database    string `mapstructure:"database" json:"database"`
	Host        string `mapstructure:"host"     json:"host"`
	Port        int    `mapstructure:"port"     json:"port"`
	Query       string `mapstructure:"query"    json:"query"`
	LogMode     bool   `mapstructure:"base-mode" json:"logMode"`
	TablePrefix string `mapstructure:"table-prefix" json:"tablePrefix"`
	Charset     string `mapstructure:"charset" json:"charset"`
	Collation   string `mapstructure:"collation" json:"collation"`
}

func NewOptions(v *viper.Viper, logger *zap.Logger) (*MysqlOptions, error) {
	var err error
	o := new(MysqlOptions)
	if err = v.UnmarshalKey("mysql", o); err != nil {
		logger.Sugar().Panicf("初始化 mysql 配置失败:%s \n", err)
		return nil, err
	}
	logger.Info("加载 mysql 配置成功")

	return o, err
}

func NewMultiDb(v *viper.Viper, dbname []string, logger *zap.Logger) (map[string]*gorm.DB, map[string]*MysqlOptions) {
	dbs := make(map[string]*gorm.DB, 0)
	dbOptions := make(map[string]*MysqlOptions, 0)
	var err error
	for _, value := range dbname {
		o := new(MysqlOptions)
		if err = v.UnmarshalKey("mysql."+value, o); err != nil {
			logger.Sugar().Panicf("初始化 mysql 配置失败:%s \n", err)
			return nil, nil
		}
		db := New(o, logger)
		if db == nil {
			logger.Sugar().Panicf("初始化 mysql 配置失败:%s \n", err)
		}
		dbs[value] = db
		dbOptions[value] = o
	}
	logger.Info("加载 mysql 配置成功")
	return dbs, dbOptions
}

// 初始化mysql数据库
func New(o *MysqlOptions, logger *zap.Logger) *gorm.DB {
	var database *gorm.DB
	tablePrefix := o.TablePrefix
	if tablePrefix != "" {
		if !strings.HasSuffix(tablePrefix, "_") {
			tablePrefix = tablePrefix + "_"
		}
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&collation=%s&%s",
		o.Username,
		o.Password,
		o.Host,
		o.Port,
		o.Database,
		o.Charset,
		o.Collation,
		o.Query,
	)
	// 隐藏密码
	showDsn := fmt.Sprintf(
		"%s:******@tcp(%s:%d)/%s?charset=%s&collation=%s&%s",
		o.Username,
		o.Host,
		o.Port,
		o.Database,
		o.Charset,
		o.Collation,
		o.Query,
	)
	//Log.Info("数据库连接DSN: ", showDsn)
	var err error
	database, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 禁用外键(指定外键时不会在mysql创建真实的外键约束)
		DisableForeignKeyConstraintWhenMigrating: true,
		//// 指定表前缀
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   tablePrefix,
			SingularTable: true, // 使用单数表名，去掉表名后缀
		},
	})
	if err != nil {
		logger.Sugar().Panicf("初始化mysql数据库异常: %v", err)
		//panic(fmt.Errorf("初始化mysql数据库异常: %v", err))
		return nil
	}

	// 开启mysql日志
	if o.LogMode {
		database.Debug()
	}

	// 全局DB赋值
	// 自动迁移表结构
	//if o.InitData {
	//	dbAutoMigrate(database, logger.Sugar())
	//}
	logger.Sugar().Infof("初始化mysql数据库完成! dsn: %s", showDsn)
	return database
}
