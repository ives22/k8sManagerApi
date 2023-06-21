package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

var DaemonSet daemonSet

type daemonSet struct{}

// DaemonSetsResp 定义列表的返回内容，Items是DaemonSet元素列表，DaemonSet是元素的数量
type DaemonSetsResp struct {
	Total int                `json:"total"`
	Items []appsv1.DaemonSet `json:"items"`
}

// DaemonSetSNp 用于返回namespace中deployment的数量
type DaemonSetSNp struct {
	Namespace    string `json:"namespace"`
	DaemonSetNum int    `json:"daemonSet_num"`
}

// DaemonSetCreate 定义DaemonSetCreate结构体，用于创建DaemonSet需要的参数属性的定义
type DaemonSetCreate struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Image         string            `json:"image"`
	Label         map[string]string `json:"label"`
	Cpu           string            `json:"cpu"`
	Memory        string            `json:"memory"`
	ContainerPort int32             `json:"container_port"`
	HealthCheck   bool              `json:"health_check"`
	HealthPath    string            `json:"health_path"`
	Cluster       string            `json:"cluster"`
}

// GetDaemonSets 获取DaemonSet列表
func (d *daemonSet) GetDaemonSets(client *kubernetes.Clientset, filterName, namespace string, limit, page int) (daemonSetRest *DaemonSetsResp, err error) {
	// 获取DaemonSetList类型的Daemonset列表
	daemonSetList, err := client.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取DaemonSet列表失败, %v", err.Error()))
		return nil, errors.New("获取DaemonSet列表失败," + err.Error())
	}
	// 实例化dataSelector结构体，组装数据
	selectableData := &dataSelector{
		GenericDataList: d.toCells(daemonSetList.Items),
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
	//将[]DataCell类型的DaemonSet列表转为appsv1.DaemonSet列表
	daemonSets := d.fromCells(data.GenericDataList)
	// 拼接返回数据
	daemonSetRest = &DaemonSetsResp{
		Total: total,
		Items: daemonSets,
	}
	return daemonSetRest, nil
}

// GetDaemonSetDetail 获取DaemonSet详情
func (d *daemonSet) GetDaemonSetDetail(client *kubernetes.Clientset, daemonSetName, namespace string) (daemonset *appsv1.DaemonSet, err error) {
	daemonset, err = client.AppsV1().DaemonSets(namespace).Get(context.TODO(), daemonSetName, metav1.GetOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取DaemonSet详情失败, %v", err.Error()))
		return nil, errors.New("获取DaemonSet详情失败, " + err.Error())
	}
	return daemonset, nil
}

// DeleteDaemonSet 删除DaemonSet
func (d *daemonSet) DeleteDaemonSet(client *kubernetes.Clientset, daemonSetName, namespace string) (err error) {
	err = client.AppsV1().DaemonSets(namespace).Delete(context.TODO(), daemonSetName, metav1.DeleteOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("删除DaemonSet失败, %v", err.Error()))
		return errors.New("删除DaemonSet失败, " + err.Error())
	}
	return nil
}

// UpdateDaemonSet 更新DaemonSet
func (d *daemonSet) UpdateDaemonSet(client *kubernetes.Clientset, namespace, content string) (err error) {
	var daemon = &appsv1.DaemonSet{}
	err = json.Unmarshal([]byte(content), daemon)
	if err != nil {
		zap.L().Error(fmt.Sprintf("反序列化失败, %v", err.Error()))
		return errors.New("反序列化失败," + err.Error())
	}
	_, err = client.AppsV1().DaemonSets(namespace).Update(context.TODO(), daemon, metav1.UpdateOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("更新DaemonSet失败, %v", err.Error()))
		return errors.New("更新DaemonSet失败, " + err.Error())
	}
	return nil
}

// CreateDaemonSet 创建DaemonSet
func (d *daemonSet) CreateDaemonSet(client *kubernetes.Clientset, data *DaemonSetCreate) (err error) {
	// 初始化一个apps
	daemonset := &appsv1.DaemonSet{
		// ObjectMeta 定义资源名、名称空间、以及标签
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		// Spec中定义daemonset的选择弃、以及Pod属性
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: data.Label,
			},
			// Pod模板信息
			Template: corev1.PodTemplateSpec{
				// Pod 元数据
				ObjectMeta: metav1.ObjectMeta{
					Name:   data.Name,
					Labels: data.Label,
				},
				// Pod spec信息，容器的名字、端口、镜像
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  data.Name,
							Image: data.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: data.ContainerPort,
									Protocol:      corev1.ProtocolTCP,
								},
							},
						},
					},
				},
			},
		},
		Status: appsv1.DaemonSetStatus{},
	}
	// 判断是否打开了检查功能
	if data.HealthCheck {
		daemonset.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					Port: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			InitialDelaySeconds: 5,
			TimeoutSeconds:      5,
			PeriodSeconds:       5,
		}
		daemonset.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					Port: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			InitialDelaySeconds: 15,
			TimeoutSeconds:      5,
			PeriodSeconds:       5,
		}
	}
	// 定义容器的limit和request资源
	daemonset.Spec.Template.Spec.Containers[0].Resources.Limits = map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:    resource.MustParse(data.Cpu),
		corev1.ResourceMemory: resource.MustParse(data.Memory),
	}
	daemonset.Spec.Template.Spec.Containers[0].Resources.Requests = map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:    resource.MustParse(data.Cpu),
		corev1.ResourceMemory: resource.MustParse(data.Memory),
	}

	// 调用sdk创建deployment
	if _, err = client.AppsV1().DaemonSets(data.Namespace).Create(context.TODO(), daemonset, metav1.CreateOptions{}); err != nil {
		zap.L().Error(fmt.Sprintf("创建DaemonSet失败, %v", err.Error()))
		return errors.New("创建DaemonSet失败," + err.Error())
	}

	return nil
}

// GetDaemonSetNumPerNp 获取每个Namespace的DaemonSet的数量
func (d *daemonSet) GetDaemonSetNumPerNp(client *kubernetes.Clientset) (DaemonSetSNps []*DaemonSetSNp, err error) {
	// 获取Namespace列表
	namespaceList, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取Namespace列表失败, %v", err.Error()))
		return nil, errors.New("获取Namespace列表失败, " + err.Error())
	}
	for _, namespace := range namespaceList.Items {
		// 获取Deployment列表
		daemonSetList, err := client.AppsV1().DaemonSets(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			zap.L().Error(fmt.Sprintf("获取DaemonSet列表失败, %v", err.Error()))
			return nil, errors.New("获取DaemonSet列表失败, " + err.Error())
		}
		// 组装数据
		daemonNp := &DaemonSetSNp{
			Namespace:    namespace.Name,
			DaemonSetNum: len(daemonSetList.Items),
		}
		// 添加数据到podsNps中
		DaemonSetSNps = append(DaemonSetSNps, daemonNp)
	}
	return DaemonSetSNps, nil
}

// 类型转换的方法，appsv1.DaemonSet-> DataCell, DataCell -> appsv1.DaemonSet
// toCells appsv1.DaemonSet -> DataCell
func (d *daemonSet) toCells(daemonSets []appsv1.DaemonSet) []DataCell {
	cells := make([]DataCell, len(daemonSets))
	for i := range daemonSets {
		cells[i] = daemonSetCell(daemonSets[i])
	}
	return cells
}

// fromCells DataCell -> appsv1.DaemonSet
func (d *daemonSet) fromCells(cells []DataCell) []appsv1.DaemonSet {
	daemonSets := make([]appsv1.DaemonSet, len(cells))
	for i := range cells {
		//  cells[i].(daemonSetCell) 是将DataCell类型转成daemonSetCell
		daemonSets[i] = appsv1.DaemonSet(cells[i].(daemonSetCell))
	}
	return daemonSets
}