package dao

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"k8sManagerApi/db/mysql"
	"k8sManagerApi/model"
)

var Workflow workflow

type workflow struct{}

type WorkflowResponse struct {
	Items []*model.Workflow `json:"items"`
	Total int32             `json:"total"`
}

// GetWorkflows 获取workflow列表
func (w *workflow) GetWorkflows(filterName, namespace, cluster string, limit, page int) (data *WorkflowResponse, err error) {
	// 定义分页的起始位置
	startSet := (page - 1) * limit
	// 定义数据库查询返回的内容
	var workflowList []*model.Workflow
	// 数据库查询，Limit方法用于限制条数，offset方法用于设置起始位置
	tx := mysql.DB.Where("name like ? and namespace = ? and cluster = ?", "%"+filterName+"%", namespace, cluster).
		Limit(limit).
		Offset(startSet).
		Order("id desc").
		Find(&workflowList)
	//gorm会默认把空数据也放到err中，故这里要排除空数据的情况
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		zap.L().Error(fmt.Sprintf("获取workflow列表失败, %v", tx.Error.Error()))
		return nil, errors.New("获取workflow列表失败," + tx.Error.Error())
	}
	data = &WorkflowResponse{
		Items: workflowList,
		Total: int32(len(workflowList)),
	}
	return data, nil
}

// GetById 查询workflow单条数据
func (w *workflow) GetById(id int) (workflow *model.Workflow, err error) {
	workflow = &model.Workflow{}
	tx := mysql.DB.Where("id =?", id).First(&workflow)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		zap.L().Error(fmt.Sprintf("获取workflow单条数据失败, %v", tx.Error.Error()))
		return nil, errors.New("获取workflow单条数据失败," + tx.Error.Error())
	}
	return workflow, nil
}

// Add 新增workflow
func (w *workflow) Add(workflow *model.Workflow) (err error) {
	tx := mysql.DB.Create(&workflow)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		zap.L().Error(fmt.Sprintf("添加workflow失败, %v", tx.Error.Error()))
		return errors.New("添加workflow失败," + tx.Error.Error())
	}
	return nil
}

// DelById 删除workflow
/*
软删除 db.GORM.Delete("id = ?", id)
软删除执行的是UPDATE语句，将deleted_at字段设置为时间即可，gorm 默认就是软删。
实际执行语句 UPDATE `workflow` SET `deleted_at` = '2021-03-01 08:32:11' WHERE `id` IN('1')
硬删除 db.GORM.Unscoped().Delete("id = ?", id)) 直接从表中删除这条数据
实际执行语句 DELETE FROM `workflow` WHERE `id` IN ('1');
*/
func (w *workflow) DelById(id int) (err error) {
	tx := mysql.DB.Where("id = ?", id).Delete(&model.Workflow{})
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		zap.L().Error(fmt.Sprintf("删除workflow失败, %v", tx.Error.Error()))
		return errors.New("删除workflow失败," + tx.Error.Error())
	}
	return nil
}