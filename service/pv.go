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

var Pv pv

type pv struct{}

// pvResp 定义列表的返回内容，Items是pv元素列表，Total是元素的数量
type pvResp struct {
	Total int                       `json:"total"`
	Items []corev1.PersistentVolume `json:"items"`
}

// GetPvs 获取pv列表
func (p *pv) GetPvs(client *kubernetes.Clientset, filterName string, limit, page int) (pvRest *pvResp, err error) {
	// 获取NodeList类型的Node列表
	pvList, err := client.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取pv列表失败, %v", err.Error()))
		return nil, errors.New("获取pv列表失败," + err.Error())
	}
	// 实例化dataSelector结构体，组装数据
	selectableData := &dataSelector{
		GenericDataList: p.toCells(pvList.Items),
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
	//将[]DataCell类型的pv列表转为corev1.PersistentVolumes列表
	pvs := p.fromCells(data.GenericDataList)
	// 拼接返回数据
	pvRest = &pvResp{
		Total: total,
		Items: pvs,
	}
	return pvRest, nil
}

// GetPvDetail 获取pv详情
func (p *pv) GetPvDetail(client *kubernetes.Clientset, pvName string) (pv *corev1.PersistentVolume, err error) {
	pv, err = client.CoreV1().PersistentVolumes().Get(context.TODO(), pvName, metav1.GetOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取pv详情失败, %v", err.Error()))
		return nil, errors.New("获取pv详情失败, " + err.Error())
	}
	return pv, nil
}

// DeletePv 删除pv
func (p *pv) DeletePv(client *kubernetes.Clientset, pvName string) (err error) {
	err = client.CoreV1().PersistentVolumes().Delete(context.TODO(), pvName, metav1.DeleteOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("删除pv失败, %v", err.Error()))
		return errors.New("删除pv失败, " + err.Error())
	}
	return nil
}

// 类型转换的方法，CoreV1.PersistentVolume-> DataCell, DataCell -> CoreV1.PersistentVolume
// toCells CoreV1.PersistentVolume -> DataCell
func (p *pv) toCells(pvs []corev1.PersistentVolume) []DataCell {
	cells := make([]DataCell, len(pvs))
	for i := range pvs {
		cells[i] = pvCell(pvs[i])
	}
	return cells
}

// fromCells DataCell -> CoreV1.PersistentVolume
func (p *pv) fromCells(cells []DataCell) []corev1.PersistentVolume {
	pvs := make([]corev1.PersistentVolume, len(cells))
	for i := range cells {
		//  cells[i].(nodeCell) 是将DataCell类型转成nodeCell
		pvs[i] = corev1.PersistentVolume(cells[i].(pvCell))
	}
	return pvs
}