package dao

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"k8sManagerApi/db/mysql"
	"k8sManagerApi/model"
	"time"
)

var Event event

type event struct{}

// Events 定义event列表的结构体
type Events struct {
	Total int64          `json:"total"`
	Items []*model.Event `json:"items"`
}

// GetEvents 获取event列表
func (e *event) GetEvents(name, cluster string, page, limit int) (events *Events, err error) {
	// 定义分页数据的起始位置
	startSet := (page - 1) * limit
	// 定义数据库查询的返回内容
	var (
		eventList       = make([]*model.Event, 0)
		total     int64 = 0
	)
	// 数据库查询
	tx := mysql.DB.Model(&model.Event{}).
		Where("name like ? and cluster = ?", "%"+name+"%", cluster).
		Count(&total).
		Limit(limit).
		Offset(startSet).
		Order("id desc").
		Find(&eventList)
	if tx.Error != nil {
		zap.L().Error(fmt.Sprintf("获取Event列表失败, %v", tx.Error))
		return nil, errors.New(fmt.Sprintf("获取Event列表失败, %v", tx.Error))
	}
	events = &Events{
		Items: eventList,
		Total: total,
	}
	return events, err
}

// Add 新增event
func (e *event) Add(event *model.Event) (err error) {
	tx := mysql.DB.Create(&event)
	if tx.Error != nil {
		zap.L().Error(fmt.Sprintf("新增Event失败, %v", tx.Error))
		return errors.New(fmt.Sprintf("新增Event失败, %v", tx.Error))
	}
	return nil
}

// HasEvent 查询单个event
func (e *event) HasEvent(name, kind, namespace, reason string, eventTime time.Time, cluster string) (*model.Event, bool, error) {
	data := &model.Event{}
	tx := mysql.DB.Where("name = ? and kind = ? and namespace = ? and reason = ? and event_time = ? and cluster = ?", name, kind, namespace, reason, eventTime, cluster).First(&data)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if tx.Error != nil {
		zap.L().Error(fmt.Sprintf("查询Event失败, %v", tx.Error))
		return nil, false, errors.New(fmt.Sprintf("查询Event失败, %v", tx.Error))
	}
	return data, true, nil
}