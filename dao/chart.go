package dao

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"k8sManagerApi/db/mysql"
	"k8sManagerApi/model"
)

var Chart chart

type chart struct{}

// Charts 定义返回内容结构体
type Charts struct {
	Items []*model.Chart `json:"items"`
	Total int64          `json:"total"`
}

// GetList 获取Chart列表
func (c *chart) GetList(name string, page, limit int) (charts *Charts, err error) {
	// 定义分页数据的起始位置
	startSet := (page - 1) * limit
	// 定义数据库查询的返回内容
	var (
		chartList []*model.Chart
		total     int64 = 0
	)
	// 数据库查询，Limit方法用于限制条数，Offset方法设置起始位置
	tx := mysql.DB.Model(&model.Chart{}).
		Where("name like ?", "%"+name+"%").
		Count(&total).
		Limit(limit).
		Offset(startSet).
		Order("id desc").
		Find(&chartList)
	if tx.Error != nil {
		zap.L().Error(fmt.Sprintf("获取Chart列表失败, %v", tx.Error))
		return nil, errors.New(fmt.Sprintf("获取Chart列表失败, %v", tx.Error))
	}
	charts = &Charts{
		Items: chartList,
		Total: total,
	}
	return charts, nil
}

// HasChart 查询单个chart
func (c *chart) HasChart(name string) (*model.Chart, bool, error) {
	data := &model.Chart{}
	tx := mysql.DB.Where("name = ?", name).First(&data)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if tx.Error != nil {
		zap.L().Error(fmt.Sprintf("查询Chart失败, %v", tx.Error))
		return nil, false, errors.New(fmt.Sprintf("查询Chart失败, %v", tx.Error))
	}
	return data, true, nil
}

// Add 新增Chart
func (c *chart) Add(chart *model.Chart) (err error) {
	tx := mysql.DB.Create(&chart)
	if tx.Error != nil {
		zap.L().Error(fmt.Sprintf("新增Chart失败, %v", tx.Error))
		return errors.New(fmt.Sprintf("新增Chart失败, %v", tx.Error))
	}
	return nil
}

// Update 更新Chart
func (c *chart) Update(chart *model.Chart) (err error) {
	tx := mysql.DB.Model(&chart).Updates(&model.Chart{
		Name:     chart.Name,
		FileName: chart.FileName,
		IconUrl:  chart.IconUrl,
		Version:  chart.Version,
		Describe: chart.Describe,
	})
	if tx.Error != nil {
		zap.L().Error(fmt.Sprintf("更新Chart失败, %v", tx.Error))
		return errors.New(fmt.Sprintf("更新Chart失败, %v", tx.Error))
	}
	return nil
}

// Delete 删除Chart
func (c *chart) Delete(id uint) (err error) {
	data := &model.Chart{}
	data.ID = uint(id)
	tx := mysql.DB.Delete(&data)
	if tx.Error != nil {
		zap.L().Error(fmt.Sprintf("删除Chart失败, %v", tx.Error))
		return errors.New(fmt.Sprintf("删除Chart失败, %v", tx.Error))
	}
	return nil
}