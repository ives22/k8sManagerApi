package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var StatefulSet statefulSet

type statefulSet struct{}

// StatefulSetsResp 定义列表的返回内容，Items是StatefulSet元素列表，StatefulSet是元素的数量
type StatefulSetsResp struct {
	Total int                  `json:"total"`
	Items []appsv1.StatefulSet `json:"items"`
}

// StatefulSetSNp 用于返回namespace中StatefulSet的数量
type StatefulSetSNp struct {
	Namespace      string `json:"namespace"`
	StatefulSetNum int    `json:"statefulset_num"`
}

// GetStatefulSets 获取StatefulSet列表
func (s *statefulSet) GetStatefulSets(client *kubernetes.Clientset, filterName, namespace string, limit, page int) (statefulSetRest *StatefulSetsResp, err error) {
	// 获取StatefulSetList类型的Statefulset列表
	statefulSetList, err := client.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("获取StatefulSet列表失败, %v\n", err.Error())
		return nil, errors.New("获取StatefulSet列表失败," + err.Error())
	}
	// 实例化dataSelector结构体，组装数据
	selectableData := &dataSelector{
		GenericDataList: s.toCells(statefulSetList.Items),
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
	//将[]DataCell类型的DStatefulSet列表转为appsv1.StatefulSet列表
	statefulSets := s.fromCells(data.GenericDataList)
	// 拼接返回数据
	statefulSetRest = &StatefulSetsResp{
		Total: total,
		Items: statefulSets,
	}
	return statefulSetRest, nil
}

// GetStatefulSetDetail 获取StatefulSet详情
func (s *statefulSet) GetStatefulSetDetail(client *kubernetes.Clientset, statefulSetName, namespace string) (statefulset *appsv1.StatefulSet, err error) {
	statefulset, err = client.AppsV1().StatefulSets(namespace).Get(context.TODO(), statefulSetName, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("获取StatefulSet详情失败, %v\n", err.Error())
		return nil, errors.New("获取StatefulSet详情失败, " + err.Error())
	}
	return statefulset, nil
}

// DeleteStatefulSet 删除StatefulSet
func (s *statefulSet) DeleteStatefulSet(client *kubernetes.Clientset, statefulSetName, namespace string) (err error) {
	err = client.AppsV1().StatefulSets(namespace).Delete(context.TODO(), statefulSetName, metav1.DeleteOptions{})
	if err != nil {
		fmt.Printf("删除StatefulSet失败, %v\n", err.Error())
		return errors.New("删除StatefulSet失败, " + err.Error())
	}
	return nil
}

// UpdateStatefulSet 更新StatefulSet
func (s *statefulSet) UpdateStatefulSet(client *kubernetes.Clientset, namespace, content string) (err error) {
	var stateful = &appsv1.StatefulSet{}
	err = json.Unmarshal([]byte(content), stateful)
	if err != nil {
		fmt.Printf("反序列化失败 %v\n", err.Error())
		return errors.New("反序列化失败," + err.Error())
	}
	_, err = client.AppsV1().StatefulSets(namespace).Update(context.TODO(), stateful, metav1.UpdateOptions{})
	if err != nil {
		fmt.Printf("更新StatefulSet失败, %v\n", err.Error())
		return errors.New("更新StatefulSet失败, " + err.Error())
	}
	return nil
}

// GetStatefulSetNumPerNp 获取每个Namespace的StatefulSet的数量
func (s *statefulSet) GetStatefulSetNumPerNp(client *kubernetes.Clientset) (StatefulSetSNps []*StatefulSetSNp, err error) {
	// 获取Namespace列表
	namespaceList, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("获取Namespace列表失败, %v\n", err.Error())
		return nil, errors.New("获取Namespace列表失败, " + err.Error())
	}
	for _, namespace := range namespaceList.Items {
		// 获取Deployment列表
		statefulSetList, err := client.AppsV1().StatefulSets(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Printf("获取StatefulSet列表失败, %v\n", err.Error())
			return nil, errors.New("获取StatefulSet列表失败, " + err.Error())
		}
		// 组装数据
		statefulNp := &StatefulSetSNp{
			Namespace:      namespace.Name,
			StatefulSetNum: len(statefulSetList.Items),
		}
		// 添加数据到podsNps中
		StatefulSetSNps = append(StatefulSetSNps, statefulNp)
	}
	return StatefulSetSNps, nil
}

// 类型转换的方法，appsv1.StatefulSet-> DataCell, DataCell -> appsv1.StatefulSet
// toCells appsv1.StatefulSet -> DataCell
func (s *statefulSet) toCells(statefulSets []appsv1.StatefulSet) []DataCell {
	cells := make([]DataCell, len(statefulSets))
	for i := range statefulSets {
		cells[i] = statefulSetCell(statefulSets[i])
	}
	return cells
}

// fromCells DataCell -> appsv1.StatefulSet
func (s *statefulSet) fromCells(cells []DataCell) []appsv1.StatefulSet {
	statefulSets := make([]appsv1.StatefulSet, len(cells))
	for i := range cells {
		//  cells[i].(StatefulSetCell) 是将DataCell类型转成StatefulSetCell
		statefulSets[i] = appsv1.StatefulSet(cells[i].(statefulSetCell))
	}
	return statefulSets
}