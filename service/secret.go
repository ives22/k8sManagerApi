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

var Secret secret

type secret struct{}

type SecretResp struct {
	Total int             `json:"total"`
	Items []corev1.Secret `json:"items"`
}

// GetSecrets 获取Secret列表
func (s *secret) GetSecrets(client *kubernetes.Clientset, namespaces, filterName string, limit, page int) (SecretRest *SecretResp, err error) {
	// 获取Secret
	SecretList, err := client.CoreV1().Secrets(namespaces).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取Secret列表失败, %v", err.Error()))
		return nil, errors.New("获取Secret列表失败," + err.Error())
	}
	// 实例化dataSelector结构体，组装数据
	selectableData := &dataSelector{
		GenericDataList: s.toCells(SecretList.Items),
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
	Secrets := s.fromCells(data.GenericDataList)
	SecretRest = &SecretResp{
		Total: total,
		Items: Secrets,
	}
	return SecretRest, nil
}

// GetSecretDetail 获取Secret详情
func (s *secret) GetSecretDetail(client *kubernetes.Clientset, namespace, SecretName string) (Secret *corev1.Secret, err error) {
	Secret, err = client.CoreV1().Secrets(namespace).Get(context.TODO(), SecretName, metav1.GetOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取Secret详情失败, %v", err.Error()))
		return nil, errors.New("获取Secret详情失败, " + err.Error())
	}
	return Secret, nil
}

// UpdateSecret 更新Secret
func (s *secret) UpdateSecret(client *kubernetes.Clientset, namespace, content string) (err error) {
	var Secret = &corev1.Secret{}
	err = json.Unmarshal([]byte(content), Secret)
	if err != nil {
		zap.L().Error(fmt.Sprintf("反序列化失败, %v", err.Error()))
		return errors.New("反序列化失败," + err.Error())
	}
	_, err = client.CoreV1().Secrets(namespace).Update(context.TODO(), Secret, metav1.UpdateOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("更新Secret失败, %v", err.Error()))
		return errors.New("更新Secret失败, " + err.Error())
	}
	return nil
}

// DeleteSecret 删除Secret
func (s *secret) DeleteSecret(client *kubernetes.Clientset, namespace, SecretName string) (err error) {
	err = client.CoreV1().Secrets(namespace).Delete(context.TODO(), SecretName, metav1.DeleteOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("删除Secret失败, %v", err.Error()))
		return errors.New("删除Secret失败, " + err.Error())
	}
	return nil
}

// 类型转换的方法，CoreV1.Namespace-> DataCell, DataCell -> CoreV1.Namespace
// toCells CoreV1.Namespace -> DataCell
func (s *secret) toCells(secrets []corev1.Secret) []DataCell {
	cells := make([]DataCell, len(secrets))
	for i := range secrets {
		cells[i] = secretCell(secrets[i])
	}
	return cells
}

// fromCells DataCell -> CoreV1.Namespace
func (s *secret) fromCells(cells []DataCell) []corev1.Secret {
	secrets := make([]corev1.Secret, len(cells))
	for i := range cells {
		//  cells[i].(nodeCell) 是将DataCell类型转成nodeCell
		secrets[i] = corev1.Secret(cells[i].(secretCell))
	}
	return secrets
}