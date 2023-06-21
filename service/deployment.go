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
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"strconv"
	"time"
)

var Deployment deployment

type deployment struct{}

// DeploymentsResp 定义列表的返回内容，Items是Deployment元素列表，Deployment是元素的数量
type DeploymentsResp struct {
	Total int                 `json:"total"`
	Items []appsv1.Deployment `json:"items"`
}

// DeployCreate 定义DeployCreate结构体，用于创建deployment需要的参数属性的定义
type DeployCreate struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Replicas      int32             `json:"replicas"`
	Image         string            `json:"image"`
	Label         map[string]string `json:"label"`
	Cpu           string            `json:"cpu"`
	Memory        string            `json:"memory"`
	ContainerPort int32             `json:"container_port"`
	HealthCheck   bool              `json:"health_check"`
	HealthPath    string            `json:"health_path"`
	Cluster       string            `json:"cluster"`
}

// DeploySNp 用于返回namespace中deployment的数量
type DeploySNp struct {
	Namespace     string `json:"namespace"`
	DeploymentNum int    `json:"deployment_num"`
}

// GetDeployments 获取Deployment列表，支持过滤、排序、分页
func (d *deployment) GetDeployments(client *kubernetes.Clientset, filterName, namespace string, limit, page int) (deploymentsRest *DeploymentsResp, err error) {
	// 获取DeploymentList类型的deployment列表
	deploymentList, err := client.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取Deployment列表失败, %v", err.Error()))
		return nil, errors.New("获取Deployment列表失败," + err.Error())
	}
	// 实例化dataSelector结构体，组装数据
	selectableData := &dataSelector{
		GenericDataList: d.toCells(deploymentList.Items),
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
	//将[]DataCell类型的deployment列表转为appsv1.deployment列表
	deployments := d.fromCells(data.GenericDataList)

	// 拼接返回数据
	deploymentsRest = &DeploymentsResp{
		Total: total,
		Items: deployments,
	}
	return deploymentsRest, nil
}

// GetDeploymentDetail 获取Deployment详情
func (d *deployment) GetDeploymentDetail(client *kubernetes.Clientset, deploymentName, namespace string) (deployment *appsv1.Deployment, err error) {
	deployment, err = client.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取Deployment详情失败, %v", err.Error()))
		return nil, errors.New("获取Deployment详情失败, " + err.Error())
	}
	return deployment, nil
}

// SetDeploymentReplicas 设置Deployment副本数
func (d *deployment) SetDeploymentReplicas(client *kubernetes.Clientset, deploymentName, namespace string, replicas int32) (replica int32, err error) {
	//获取autoscalingv1.Scale类型的对象，能点出当前的副本数
	scale, err := client.AppsV1().Deployments(namespace).GetScale(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取Deployment副本数失败, %v", err.Error()))
		return 0, errors.New("获取Deployment副本数失败," + err.Error())
	}
	// 修改副本数
	scale.Spec.Replicas = replicas
	_, err = client.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), deploymentName, scale, metav1.UpdateOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("更新Deployment副本数失败, %v", err.Error()))
		return 0, errors.New("更新Deployment副本数失败," + err.Error())
	}
	return scale.Spec.Replicas, nil
}

// CreateDeployment 创建Deployment，接收DeployCreate对象
func (d *deployment) CreateDeployment(client *kubernetes.Clientset, data *DeployCreate) (err error) {
	// 1、初始化一个appsv1.Deployment类型的对象，并将入参的data数据放进去
	deployment := &appsv1.Deployment{
		// ObjectMeta中定义资源名、名称空间、以及标签
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		// Spec中定义deploy的副本数，选择器，以及Pod属性
		Spec: appsv1.DeploymentSpec{
			Replicas: &data.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: data.Label,
			},
			// Pod 模版信息
			Template: corev1.PodTemplateSpec{
				// Pod 元信息，名称、标签
				ObjectMeta: metav1.ObjectMeta{
					Name:   data.Name,
					Labels: data.Label,
				},
				// Pod spec信息 定义容器名、镜像和端口
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  data.Name,
							Image: data.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: data.ContainerPort,
								},
							},
						},
					},
				},
			},
		},
		//Status定义资源的运行状态，这里由于是新建，传入空的appsv1.DeploymentStatus{}对象即可
		Status: appsv1.DeploymentStatus{},
	}
	// 2、判断检查功能是否打开，若打开，则增加健康检查功能  (这里需要优化，这里只是给第一个容器加了健康检查)
	if data.HealthCheck {
		deployment.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					//intstr.IntOrString的作用是端口可以定义为整型，也可以定义为字符串
					//Type=0则表示表示该结构体实例内的数据为整型，转json时只使用IntVal的数据
					//Type=1则表示表示该结构体实例内的数据为字符串，转json时只使用StrVal的数据
					Port: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			// 初始化等待时间
			InitialDelaySeconds: 5,
			// 超时时间
			TimeoutSeconds: 5,
			// 执行间隔
			PeriodSeconds: 5,
		}
		deployment.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
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
	deployment.Spec.Template.Spec.Containers[0].Resources.Limits = map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:    resource.MustParse(data.Cpu),
		corev1.ResourceMemory: resource.MustParse(data.Memory),
	}
	deployment.Spec.Template.Spec.Containers[0].Resources.Requests = map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:    resource.MustParse(data.Cpu),
		corev1.ResourceMemory: resource.MustParse(data.Memory),
	}

	// 调用sdk创建deployment
	_, err = client.AppsV1().Deployments(data.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("创建Deployment失败, %v", err.Error()))
		return errors.New("创建Deployment失败," + err.Error())
	}
	return nil
}

// DeleteDeployment 删除Deployment
func (d *deployment) DeleteDeployment(client *kubernetes.Clientset, deploymentName, namespace string) (err error) {
	err = client.AppsV1().Deployments(namespace).Delete(context.TODO(), deploymentName, metav1.DeleteOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("删除Deployment失败, %v", err.Error()))
		return errors.New("删除Deployment失败," + err.Error())
	}
	return nil
}

// RestartDeployment 重启Deployment
func (d *deployment) RestartDeployment(client *kubernetes.Clientset, deploymentName, namespace string) (err error) {
	//此功能等同于一下kubectl命令
	// kubectl patch deployment my-deployment -p \
	//'{"spec":{"template":{"spec":{"containers":[{"name":"my-container","env":[{"name":"MY_ENV","value":"my-value"}]}]}}}}'
	// 首先获取该deployment中每个Pod有多少容器，拿到名字
	deploymentInfo, err := client.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取Deployment失败, %v", err.Error()))
		return errors.New("获取Deployment失败," + err.Error())
	}
	// 遍历deploymentInfo获取所有容器名，放入到containerList
	var containerListData []map[string]interface{}
	for _, container := range deploymentInfo.Spec.Template.Spec.Containers {
		containerData := map[string]interface{}{
			"name": container.Name,
			"env": []map[string]interface{}{
				{
					"name":  "RESTART_",
					"value": strconv.FormatInt(time.Now().Unix(), 10),
				},
			},
		}
		containerListData = append(containerListData, containerData)
	}
	// 使用patchData Map组装数据
	patchData := map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": containerListData,
				},
			},
		},
	}
	// 序列化为字节，因为patch方法只接收字节类型参数
	patchByte, err := json.Marshal(patchData)
	fmt.Println(string(patchByte))
	if err != nil {
		zap.L().Error(fmt.Sprintf("Json序列化失败, %v", err.Error()))
		return errors.New("Json序列化失败," + err.Error())
	}
	// 调用patch方法更新deployment
	_, err = client.AppsV1().Deployments(namespace).Patch(context.TODO(), deploymentName, types.StrategicMergePatchType, patchByte, metav1.PatchOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("重启Deployment失败, %v", err.Error()))
		return errors.New("重启Deployment失败," + err.Error())
	}
	return nil
}

// UpdateDeployment 更新Deployment
func (d *deployment) UpdateDeployment(client *kubernetes.Clientset, namespace, content string) (err error) {
	var deploy = &appsv1.Deployment{}
	err = json.Unmarshal([]byte(content), deploy)
	if err != nil {
		zap.L().Error(fmt.Sprintf("反序列化失败, %v", err.Error()))
		return errors.New("反序列化失败," + err.Error())
	}
	_, err = client.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("更新Deployment失败, %v", err.Error()))
		return errors.New("更新Deployment失败," + err.Error())
	}
	return nil
}

// GetDeployNumPerNp 获取每个Namespace的Deployment的数量
func (d *deployment) GetDeployNumPerNp(client *kubernetes.Clientset) (deploysNps []*DeploySNp, err error) {
	// 获取Namespace列表
	namespaceList, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取Namespace列表失败, %v", err.Error()))
		return nil, errors.New("获取Namespace列表失败, " + err.Error())
	}
	for _, namespace := range namespaceList.Items {
		// 获取Deployment列表
		deployList, err := client.AppsV1().Deployments(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			zap.L().Error(fmt.Sprintf("获取Deployment列表失败, %v", err.Error()))
			return nil, errors.New("获取Deployment列表失败, " + err.Error())
		}
		// 组装数据
		deploysNp := &DeploySNp{
			Namespace:     namespace.Name,
			DeploymentNum: len(deployList.Items),
		}
		// 添加数据到podsNps中
		deploysNps = append(deploysNps, deploysNp)
	}
	return deploysNps, nil
}

// 类型转换的方法，appsv1.Deployment -> DataCell, DataCell -> appsv1.Deployment
// toCells appsv1.Deployment -> DataCell
func (d *deployment) toCells(deployments []appsv1.Deployment) []DataCell {
	cells := make([]DataCell, len(deployments))
	for i := range deployments {
		cells[i] = deploymentCell(deployments[i])
	}
	return cells
}

// fromCells DataCell -> appsv1.Deployment
func (d *deployment) fromCells(cells []DataCell) []appsv1.Deployment {
	deployments := make([]appsv1.Deployment, len(cells))
	for i := range cells {
		//  cells[i].(deploymentCell) 是将DataCell类型转成deploymentCell
		deployments[i] = appsv1.Deployment(cells[i].(deploymentCell))
	}
	return deployments
}