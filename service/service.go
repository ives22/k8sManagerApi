package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

var Service service

type service struct{}

// ServicesResp 定义列表的返回内容，Items是Sservice元素列表，Total是元素的数量
type ServicesResp struct {
	Total int              `json:"total"`
	Items []corev1.Service `json:"items"`
}

/*
// ServiceCreate 用于创建service需要的参数属性的定义
type ServiceCreate struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Type          string            `json:"type"` // 分为ClusterIP、NodePort、LoadBalancer
	Port          int32             `json:"port"`
	ContainerPort int32             `json:"container_port"`
	NodePort      int32             `json:"node_port"`
	Label         map[string]string `json:"label"`
}
*/

// ServiceCreate 用于创建service需要的参数属性的定义
type ServiceCreate struct {
	Name           string            `json:"name"`             // service 名字
	Namespace      string            `json:"namespace"`        // 所属名称空间
	Type           string            `json:"type"`             // 分为ClusterIP、NodePort、LoadBalancer
	Selector       map[string]string `json:"selector"`         // 标签选择器
	Label          map[string]string `json:"label"`            // service自身的标签
	Port           int32             `json:"port"`             // service端口
	PortName       string            `json:"port_name"`        // 端口名称
	Protocol       string            `json:"protocol"`         // 协议，TCP、UDP、默认为TCP
	TargetPort     int32             `json:"target_port"`      // 目标端口，可以理解为容器端口
	NodePort       int32             `json:"node_port"`        // 如果类型为NodePort类型，则需要该字断
	LoadBalancerIP string            `json:"load_balancer_ip"` // 如果类型为LoadBalancer，作为可选项
	ExternalIPs    []string          `json:"externalIPs"`      // 可选项
	Cluster        string            `json:"cluster"`
}

// ServiceSNp 用于返回namespace中service的数量
type ServiceSNp struct {
	Namespace  string `json:"namespace"`
	ServiceNum int    `json:"service_num"`
}

// GetServices 获取Service列表
func (s *service) GetServices(client *kubernetes.Clientset, filterName, namespace string, limit, page int) (serviceRest *ServicesResp, err error) {
	// 获取serviceList类型的service列表
	serviceList, err := client.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取Service列表失败, %v", err.Error()))
		return nil, errors.New("获取Service列表失败," + err.Error())
	}
	// 实例化dataSelector结构体，组装数据
	selectableData := &dataSelector{
		GenericDataList: s.toCells(serviceList.Items),
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
	services := s.fromCells(data.GenericDataList)
	// 拼接返回数据
	serviceRest = &ServicesResp{
		Total: total,
		Items: services,
	}
	return serviceRest, nil
}

// GetServicesDetail 获取Service详情
func (s *service) GetServicesDetail(client *kubernetes.Clientset, serviceName, namespace string) (service *corev1.Service, err error) {
	service, err = client.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取Service详情失败, %v", err.Error()))
		return nil, errors.New("获取Service详情失败, " + err.Error())
	}
	return service, nil
}

// CreateService 创建Service，接收ServiceCreate对象
func (s *service) CreateService(client *kubernetes.Clientset, data *ServiceCreate) (err error) {
	// 将data中的数据组装为corv1.Service对象
	serviced := &corev1.Service{
		// ObjectMeta中定义资源名、名称空间以及标签
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		// Spec中定义类型，端口，选择器
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceType(data.Type),
			Ports: []corev1.ServicePort{
				{
					Name:     data.PortName,
					Port:     data.Port,
					Protocol: corev1.Protocol(data.Protocol),
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: data.TargetPort,
					},
				},
			},
			Selector: data.Selector,
		},
	}
	// 判断是否存在ExternalIPs, 如果存在则添加上
	if len(data.ExternalIPs) > 0 {
		serviced.Spec.ExternalIPs = data.ExternalIPs
	}
	// 判断类型，默认创建为ClusterIP类型，这里如果是NodePort，则添加配置NodePort的配置
	if data.Type == "NodePort" && data.NodePort != 0 {
		serviced.Spec.Ports[0].NodePort = data.NodePort
	}
	// 将service对象组装为corv1.Service对象
	_, err = client.CoreV1().Services(data.Namespace).Create(context.TODO(), serviced, metav1.CreateOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("创建Service失败, %v", err.Error()))
		return errors.New("创建Service失败, " + err.Error())
	}
	return nil
}

// DeleteService 删除Service
func (s *service) DeleteService(client *kubernetes.Clientset, serviceName, namespace string) (err error) {
	err = client.CoreV1().Services(namespace).Delete(context.TODO(), serviceName, metav1.DeleteOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("删除Service失败, %v", err.Error()))
		return errors.New("删除Service失败, " + err.Error())
	}
	return nil
}

// UpdateService 更新Service
func (s *service) UpdateService(client *kubernetes.Clientset, namespace, content string) (err error) {
	var service = &corev1.Service{}
	err = json.Unmarshal([]byte(content), service)
	if err != nil {
		zap.L().Error(fmt.Sprintf("反序列化失败, %v", err.Error()))
		return errors.New("反序列化失败," + err.Error())
	}
	_, err = client.CoreV1().Services(namespace).Update(context.TODO(), service, metav1.UpdateOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("更新Service失败, %v", err.Error()))
		return errors.New("更新Service失败, " + err.Error())
	}
	return nil
}

// 类型转换的方法，corev1.Service-> DataCell, DataCell -> corev1.Service
// toCells corev1.Service-> DataCell
func (s *service) toCells(services []corev1.Service) []DataCell {
	cells := make([]DataCell, len(services))
	for i := range services {
		cells[i] = serviceCell(services[i])
	}
	return cells
}

// fromCells DataCell -> corev1.Service
func (s *service) fromCells(cells []DataCell) []corev1.Service {
	services := make([]corev1.Service, len(cells))
	for i := range cells {
		//  cells[i].(daemonSetCell) 是将DataCell类型转成daemonSetCell
		services[i] = corev1.Service(cells[i].(serviceCell))
	}
	return services
}