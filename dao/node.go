package dao

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"k8sManagerApi/db/mysql"
	"k8sManagerApi/model"
)

var Node node

type node struct{}

// Nodes 定义Node列表的结构体，返回对应集群的所有node
type Nodes struct {
	Total int64         `json:"total"`
	Items []*model.Node `json:"items"`
}

type ClusterNodes struct {
	Cluster  string `json:"cluster"`
	IP       string `json:"ip"`
	HostName string `json:"host_name"`
}

// GetNodes 获取node列表，这里用于接口返回
func (n *node) GetNodes(cluster string, page, limit int) (nodes *Nodes, err error) {
	// 定义分页数据的起始位置
	startSet := (page - 1) * limit
	// 定义数据库查询返回的内容
	var (
		nodeList       = make([]*model.Node, 0)
		total    int64 = 0
	)
	// 数据库查询
	tx := mysql.DB.Model(&model.Node{}).
		Where("cluster = ?", cluster).
		Count(&total).
		Limit(limit).
		Offset(startSet).
		Order("id desc").
		Find(&nodeList)
	if tx.Error != nil {
		zap.L().Error(fmt.Sprintf("获取node列表失败, %v", tx.Error))
		return nil, errors.New(fmt.Sprintf("获取node列表失败, %v", tx.Error))
	}
	nodes = &Nodes{
		Items: nodeList,
		Total: total,
	}
	return nodes, nil
}

// GetNodeList 获取指定集群的node的IP信息，用于安装chart时候，传输chart文件使用。
func (n *node) GetNodeList(cluster string) (nodeList []*ClusterNodes, err error) {
	tx := mysql.DB.Debug().Model(&model.Node{}).Select("cluster", "ip", "host_name").Where("cluster = ?", cluster).Find(&nodeList)
	if tx.Error != nil {
		zap.L().Error(fmt.Sprintf("获取node列表失败, %v", tx.Error))
		return nil, errors.New(fmt.Sprintf("获取node列表失败, %v", tx.Error))
	}
	return nodeList, nil
}

// Add 新增node，当获取集群的node时候进行入库
func (n *node) Add(node *model.Node) (err error) {
	// 先判断节点是否存在，如果存在则不添加
	existingNode := &model.Node{}
	result := mysql.DB.Where("host_name = ? and cluster = ? and ip = ?", node.HostName, node.Cluster, node.IP).First(&existingNode)
	if result == nil {
		// 节点已存在，不进行插入操作
		zap.L().Warn(fmt.Sprintf("节点已存在, cluster: %v, ip: %v", node.Cluster, node.IP))
		return nil
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 查询出错
		zap.L().Error(fmt.Sprintf("查询节点失败: %v", result.Error))
		return errors.New(fmt.Sprintf("查询节点失败: %v", result.Error))
	}

	// 插入新节点
	tx := mysql.DB.Create(&node)
	if tx.Error != nil {
		zap.L().Error(fmt.Sprintf("新增node失败, %v", tx.Error))
		return errors.New(fmt.Sprintf("新增node失败, %v", tx.Error))
	}
	return nil
}