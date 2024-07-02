package dao

import (
	"gitee.com/go-nianxi/go-auth/internal/auth/model"
	"gitee.com/go-nianxi/go-common/pkg/valid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type LoginAuditLogDao struct {
	db     *gorm.DB
	valid  *valid.Validator
	logger *zap.Logger
}

func NewLoginAuditLogDao(db *gorm.DB) *LoginAuditLogDao {
	return &LoginAuditLogDao{
		db: db,
	}
}

// var Logs []entity.AuthVisitLog //全局变量多个线程需要加锁，所以每个线程自己维护一个
// 处理AuthVisitLogChan将日志记录到数据库
func (dao *LoginAuditLogDao) SaveLoginLogChannel(olc <-chan *model.LoginAuditLog) {
	// 只会在线程开启的时候执行一次
	Logs := make([]model.LoginAuditLog, 0)

	// 一直执行--收到olc就会执行
	for log := range olc {
		Logs = append(Logs, *log)
		// 每10条记录到数据库
		if len(Logs) > 5 {
			dao.db.Create(&Logs)
			Logs = make([]model.LoginAuditLog, 0)
		}
	}
}
