package service

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var Node node

type node struct{}

// nodeResp 定义列表的返回内容，Items是Node元素列表，Total是元素的数量
type nodeResp struct {
	Total int           `json:"total"`
	Items []corev1.Node `json:"items"`
}

// GetNodes 获取node列表
func (n *node) GetNodes(client *kubernetes.Clientset, filterName string, limit, page int) (nodeRest *nodeResp, err error) {
	// 获取NodeList类型的Node列表
	nodeList, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取node列表失败, %v", err.Error()))
		return nil, errors.New("获取node列表失败," + err.Error())
	}
	// 实例化dataSelector结构体，组装数据
	selectableData := &dataSelector{
		GenericDataList: n.toCells(nodeList.Items),
		DataSelect: &DataSelectQuery{
			Filter: &FilterQuery{Name: filterName},
			Paginate: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}
	// 先过滤
	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)
	// 再排序和分页
	data := filtered.Sort().Paginate()
	//将[]DataCell类型的Node列表转为corev1.node列表
	nodes := n.fromCells(data.GenericDataList)
	// 拼接返回数据
	nodeRest = &nodeResp{
		Total: total,
		Items: nodes,
	}
	return nodeRest, nil
}

// GetNodeDetail 获取node详情
func (n *node) GetNodeDetail(client *kubernetes.Clientset, nodeName string) (node *corev1.Node, err error) {
	node, err = client.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取Node详情失败, %v", err.Error()))
		return nil, errors.New("获取Node详情失败, " + err.Error())
	}

	return node, nil
}

// 类型转换的方法，CoreV1.Node-> DataCell, DataCell -> CoreV1.Node
// toCells CoreV1.Node -> DataCell
func (n *node) toCells(nodeS []corev1.Node) []DataCell {
	cells := make([]DataCell, len(nodeS))
	for i := range nodeS {
		cells[i] = nodeCell(nodeS[i])
	}
	return cells
}

// fromCells DataCell -> CoreV1.Node
func (n *node) fromCells(cells []DataCell) []corev1.Node {
	nodes := make([]corev1.Node, len(cells))
	for i := range cells {
		//  cells[i].(nodeCell) 是将DataCell类型转成nodeCell
		nodes[i] = corev1.Node(cells[i].(nodeCell))
	}
	return nodes
}