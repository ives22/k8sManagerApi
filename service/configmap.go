package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var ConfigMap configMap

type configMap struct{}

type configMapResp struct {
	Total int                `json:"total"`
	Items []corev1.ConfigMap `json:"items"`
}

// GetConfigMaps 获取configmap列表
func (c *configMap) GetConfigMaps(client *kubernetes.Clientset, namespaces, filterName string, limit, page int) (configMapRest *configMapResp, err error) {
	// 获取configmap
	configMapList, err := client.CoreV1().ConfigMaps(namespaces).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("获取configMap列表失败, %v\n", err.Error())
		return nil, errors.New("获取configMap列表失败," + err.Error())
	}
	// 实例化dataSelector结构体，组装数据
	selectableData := &dataSelector{
		GenericDataList: c.toCells(configMapList.Items),
		DataSelect: &DataSelectQuery{
			Filter: &FilterQuery{Name: filterName},
			Paginate: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}
	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)
	data := filtered.Sort().Paginate()
	configmaps := c.fromCells(data.GenericDataList)
	configMapRest = &configMapResp{
		Total: total,
		Items: configmaps,
	}
	return configMapRest, nil
}

// GetConfigMapDetail 获取configMap详情
func (c *configMap) GetConfigMapDetail(client *kubernetes.Clientset, namespace, configmapName string) (configMap *corev1.ConfigMap, err error) {
	configMap, err = client.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configmapName, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("获取ConfigMap详情失败, %v\n", err.Error())
		return nil, errors.New("获取ConfigMap详情失败, " + err.Error())
	}
	return configMap, nil
}

// UpdateConfigMap 更新ConfigMap
func (c *configMap) UpdateConfigMap(client *kubernetes.Clientset, namespace, content string) (err error) {
	var configMap = &corev1.ConfigMap{}
	err = json.Unmarshal([]byte(content), configMap)
	if err != nil {
		fmt.Printf("反序列化失败 %v\n", err.Error())
		return errors.New("反序列化失败," + err.Error())
	}
	_, err = client.CoreV1().ConfigMaps(namespace).Update(context.TODO(), configMap, metav1.UpdateOptions{})
	if err != nil {
		fmt.Printf("更新ConfigMap失败, %v\n", err.Error())
		return errors.New("更新ConfigMap失败, " + err.Error())
	}
	return nil
}

// DeleteConfigMap 删除ConfigMap
func (c *configMap) DeleteConfigMap(client *kubernetes.Clientset, namespace, configmapName string) (err error) {
	err = client.CoreV1().ConfigMaps(namespace).Delete(context.TODO(), configmapName, metav1.DeleteOptions{})
	if err != nil {
		fmt.Printf("删除ConfigMap失败, %v\n", err.Error())
		return errors.New("删除ConfigMap失败, " + err.Error())
	}
	return nil
}

// 类型转换的方法，CoreV1.Namespace-> DataCell, DataCell -> CoreV1.Namespace
// toCells CoreV1.Namespace -> DataCell
func (c *configMap) toCells(configMaps []corev1.ConfigMap) []DataCell {
	cells := make([]DataCell, len(configMaps))
	for i := range configMaps {
		cells[i] = configmapCell(configMaps[i])
	}
	return cells
}

// fromCells DataCell -> CoreV1.Namespace
func (c *configMap) fromCells(cells []DataCell) []corev1.ConfigMap {
	configMaps := make([]corev1.ConfigMap, len(cells))
	for i := range cells {
		//  cells[i].(nodeCell) 是将DataCell类型转成nodeCell
		configMaps[i] = corev1.ConfigMap(cells[i].(configmapCell))
	}
	return configMaps
}