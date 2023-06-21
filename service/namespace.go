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

var Namespace namespace

type namespace struct{}

// namespaceResp 定义列表的返回内容，Items是namespace元素列表，Total是元素的数量
type namespaceResp struct {
	Total int                `json:"total"`
	Items []corev1.Namespace `json:"items"`
}

type NamespaceCreate struct {
	Name    string            `json:"name"`
	Label   map[string]string `json:"label"`
	Cluster string            `json:"cluster"`
}

// GetNamespaces 获取Namespace列表
func (n *namespace) GetNamespaces(client *kubernetes.Clientset, filterName string, limit, page int) (namespaceRest *namespaceResp, err error) {
	// 获取NodeList类型的Node列表
	namespaceList, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取namespace列表失败, %v", err.Error()))
		return nil, errors.New("获取namespace列表失败," + err.Error())
	}
	// 实例化dataSelector结构体，组装数据
	selectableData := &dataSelector{
		GenericDataList: n.toCells(namespaceList.Items),
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
	//将[]DataCell类型的namespace列表转为corev1.namespace列表
	namespaces := n.fromCells(data.GenericDataList)
	// 拼接返回数据
	namespaceRest = &namespaceResp{
		Total: total,
		Items: namespaces,
	}
	return namespaceRest, nil
}

// GetNamespaceDetail 获取namespace详情
func (n *namespace) GetNamespaceDetail(client *kubernetes.Clientset, namespaceName string) (namespace *corev1.Namespace, err error) {
	namespace, err = client.CoreV1().Namespaces().Get(context.TODO(), namespaceName, metav1.GetOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取Namespace详情失败, %v", err.Error()))
		return nil, errors.New("获取Namespace详情失败, " + err.Error())
	}
	return namespace, nil
}

// CreateNamespace 创建namespace
func (n *namespace) CreateNamespace(client *kubernetes.Clientset, data *NamespaceCreate) (err error) {
	// 1、初始化一个appsv1.Deployment类型的对象，并将入参的data数据放进去
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   data.Name,
			Labels: data.Label,
		},
		Spec:   corev1.NamespaceSpec{},
		Status: corev1.NamespaceStatus{},
	}
	// 2、调用SDK创建namespace
	_, err = client.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("创建namespace失败, %v", err.Error()))
		return errors.New("创建namespace失败," + err.Error())
	}
	return nil
}

// DeleteNamespace 删除namespace
func (n *namespace) DeleteNamespace(client *kubernetes.Clientset, namespaceName string) (err error) {
	err = client.CoreV1().Namespaces().Delete(context.TODO(), namespaceName, metav1.DeleteOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("删除Namespace失败, %v", err.Error()))
		return errors.New("删除Namespace失败, " + err.Error())
	}
	return nil
}

// 类型转换的方法，CoreV1.Namespace-> DataCell, DataCell -> CoreV1.Namespace
// toCells CoreV1.Namespace -> DataCell
func (n *namespace) toCells(namespaces []corev1.Namespace) []DataCell {
	cells := make([]DataCell, len(namespaces))
	for i := range namespaces {
		cells[i] = namespaceCell(namespaces[i])
	}
	return cells
}

// fromCells DataCell -> CoreV1.Namespace
func (n *namespace) fromCells(cells []DataCell) []corev1.Namespace {
	namespaces := make([]corev1.Namespace, len(cells))
	for i := range cells {
		//  cells[i].(nodeCell) 是将DataCell类型转成nodeCell
		namespaces[i] = corev1.Namespace(cells[i].(namespaceCell))
	}
	return namespaces
}