package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var Pvc pvc

type pvc struct{}

type PvcResp struct {
	Total int                            `json:"total"`
	Items []corev1.PersistentVolumeClaim `json:"items"`
}

// GetPvcs 获取PVC列表
func (p *pvc) GetPvcs(client *kubernetes.Clientset, namespace, filterName string, limit, page int) (PvcRest *PvcResp, err error) {
	// 获取PVC list
	pvcList, err := client.CoreV1().PersistentVolumeClaims(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取PVC列表失败, %v", err.Error()))
		return nil, errors.New("获取PVC列表失败" + err.Error())
	}
	// 实例化dataSelector结构体，组装数据
	selectableData := &dataSelector{
		GenericDataList: p.toCells(pvcList.Items),
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
	pvcs := p.fromCells(data.GenericDataList)
	PvcRest = &PvcResp{
		Total: total,
		Items: pvcs,
	}
	return PvcRest, err
}

// GetPvcDetail 获取PVC详情
func (p *pvc) GetPvcDetail(client *kubernetes.Clientset, namespace, pvcName string) (Pvc *corev1.PersistentVolumeClaim, err error) {
	// 获取PVC详情
	Pvc, err = client.CoreV1().PersistentVolumeClaims(namespace).Get(context.TODO(), pvcName, metav1.GetOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取PVC详情失败, %v", err.Error()))
		return nil, errors.New("获取PVC详情失败" + err.Error())
	}
	return Pvc, err
}

// UpdatePvc 更新PVC
func (p *pvc) UpdatePvc(client *kubernetes.Clientset, namespace, content string) (err error) {
	var pvc = &corev1.PersistentVolumeClaim{}
	err = json.Unmarshal([]byte(content), pvc)
	if err != nil {
		zap.L().Error(fmt.Sprintf("反序列化失败, %v", err.Error()))
		return errors.New("反序列化失败," + err.Error())
	}
	_, err = client.CoreV1().PersistentVolumeClaims(namespace).Update(context.TODO(), pvc, metav1.UpdateOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("更新PVC失败, %v", err.Error()))
		return errors.New("更新PVC失败" + err.Error())
	}
	return nil
}

// DeletePvc 删除PVC
func (p *pvc) DeletePvc(client *kubernetes.Clientset, namespace, pvcName string) (err error) {
	err = client.CoreV1().PersistentVolumeClaims(namespace).Delete(context.TODO(), pvcName, metav1.DeleteOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("删除PVC失败, %v", err.Error()))
		return errors.New("删除PVC失败" + err.Error())
	}
	return nil
}

// 类型转换的方法，CoreV1.PersistentVolumeClaim-> DataCell, DataCell -> CoreV1.PersistentVolumeClaim
// toCells CoreV1.PersistentVolumeClaim -> DataCell
func (p *pvc) toCells(pvc []corev1.PersistentVolumeClaim) []DataCell {
	cells := make([]DataCell, len(pvc))
	for i := range pvc {
		cells[i] = pvcCell(pvc[i])
	}
	return cells
}

// fromCells DataCell -> CoreV1.PersistentVolumeClaim
func (p *pvc) fromCells(cells []DataCell) []corev1.PersistentVolumeClaim {
	pvc := make([]corev1.PersistentVolumeClaim, len(cells))
	for i := range cells {
		//  cells[i].(PersistentVolumeClaim) 是将DataCell类型转成PersistentVolumeClaim
		pvc[i] = corev1.PersistentVolumeClaim(cells[i].(pvcCell))
	}
	return pvc
}