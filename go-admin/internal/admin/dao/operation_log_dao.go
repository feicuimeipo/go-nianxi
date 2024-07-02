package dao

import (
	"fmt"
	"gitee.com/go-nianxi/go-admin/internal/admin/domain/vo"
	"gitee.com/go-nianxi/go-admin/internal/admin/model"
	"gorm.io/gorm"
	"strings"
)

type OperationLogDao struct {
	db *gorm.DB
}

func NewOperationLogDao(db *gorm.DB) *OperationLogDao {
	return &OperationLogDao{
		db: db,
	}
}

func (o OperationLogDao) GetOperationLogs(req *vo.OperationLogListRequest) ([]model.OperationLog, int64, error) {
	var list []model.OperationLog
	db := o.db.Model(&model.OperationLog{}).Order("start_time DESC")

	username := strings.TrimSpace(req.Username)
	if username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	ip := strings.TrimSpace(req.Ip)
	if ip != "" {
		db = db.Where("ip LIKE ?", fmt.Sprintf("%%%s%%", ip))
	}
	path := strings.TrimSpace(req.Path)
	if path != "" {
		db = db.Where("path LIKE ?", fmt.Sprintf("%%%s%%", path))
	}
	status := req.Status
	if status != 0 {
		db = db.Where("status = ?", status)
	}

	// 分页
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := req.PageNum
	pageSize := req.PageSize
	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}

	return list, total, err

}

func (o OperationLogDao) BatchDeleteOperationLogByIds(ids []uint) error {
	err := o.db.Where("id IN (?)", ids).Unscoped().Delete(&model.OperationLog{}).Error
	return err
}

// var Logs []model.OperationLog //全局变量多个线程需要加锁，所以每个线程自己维护一个
// 处理OperationLogChan将日志记录到数据库
func (o OperationLogDao) SaveOperationLogChannel(olc <-chan *model.OperationLog) {
	// 只会在线程开启的时候执行一次
	Logs := make([]model.OperationLog, 0)

	// 一直执行--收到olc就会执行
	for log := range olc {
		Logs = append(Logs, *log)
		// 每10条记录到数据库
		if len(Logs) > 5 {
			o.db.Create(&Logs)
			Logs = make([]model.OperationLog, 0)
		}
	}
}
