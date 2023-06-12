package service

import (
	"k8s.io/client-go/kubernetes"
	"k8sManagerApi/dao"
	"k8sManagerApi/model"
)

var Workflow workflow

type workflow struct{}

// WorkflowCreate 定义workflowCreate类型
type WorkflowCreate struct {
	Name          string                 `json:"name"`
	Namespace     string                 `json:"namespace"`
	Replicas      int32                  `json:"replicas"`
	Image         string                 `json:"image"`
	Label         map[string]string      `json:"label"`
	Cpu           string                 `json:"cpu"`
	Memory        string                 `json:"memory"`
	ContainerPort int32                  `json:"container_port"`
	HealthCheck   bool                   `json:"health_check"`
	HealthPath    string                 `json:"health_path"`
	Type          string                 `json:"type"`
	Port          int32                  `json:"port"`
	NodePort      int32                  `json:"node_port"`
	Hosts         map[string][]*HttpPath `json:"hosts"`
	Cluster       string                 `json:"cluster"`
}

// GetWorkflows 获取列表分页查询
func (w *workflow) GetWorkflows(filterName, namespace, cluster string, limit, page int) (data *dao.WorkflowResponse, err error) {
	data, err = dao.Workflow.GetWorkflows(filterName, namespace, cluster, limit, page)
	if err != nil {
		return nil, err
	}
	return data, err
}

// GetWorkflowDetail 查询workflow单条数据
func (w *workflow) GetWorkflowDetail(id int) (data *model.Workflow, err error) {
	data, err = dao.Workflow.GetById(id)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// CreateWorkflow 创建workflow
func (w *workflow) CreateWorkflow(client *kubernetes.Clientset, data *WorkflowCreate) (err error) {
	// 定义ingress名字
	var ingressName string
	if data.Type == "Ingress" {
		ingressName = getIngressName(data.Name)
	} else {
		ingressName = ""
	}
	// workflow数据落库
	newWorkflow := &model.Workflow{
		Name:       data.Name,
		Namespace:  data.Namespace,
		Replicas:   data.Replicas,
		Deployment: data.Name,
		Service:    getServiceName(data.Name),
		Ingress:    ingressName,
		Type:       data.Type,
		Cluster:    data.Cluster,
	}
	err = dao.Workflow.Add(newWorkflow)
	if err != nil {
		return err
	}

	// 创建k8s资源
	err = createWorkFlowResource(client, data)
	if err != nil {
		return err
	}
	return err
}

// DeleteWorkflow 删除workflow
func (w *workflow) DeleteWorkflow(client *kubernetes.Clientset, id int) (err error) {
	// 获取数据库workflow数据
	workflow, err := dao.Workflow.GetById(id)
	if err != nil {
		return err
	}
	// 删除k8资源
	if err := deleteWorkFlowResource(client, workflow); err != nil {
		return err
	}
	// 删除数据库资源
	if err := dao.Workflow.DelById(id); err != nil {
		return err
	}
	return err
}

// 创建k8s资源 deployment service ingress
func createWorkFlowResource(client *kubernetes.Clientset, data *WorkflowCreate) (err error) {
	// 组装DeployCreate类型的数据
	dc := &DeployCreate{
		Name:          data.Name,
		Namespace:     data.Namespace,
		Replicas:      data.Replicas,
		Image:         data.Image,
		Label:         data.Label,
		Cpu:           data.Cpu,
		Memory:        data.Memory,
		ContainerPort: data.Port,
		HealthCheck:   data.HealthCheck,
		HealthPath:    data.HealthPath,
	}
	// 创建deployment
	err = Deployment.CreateDeployment(client, dc)
	if err != nil {
		return err
	}

	//声明service类型
	var serviceType string
	if data.Type != "Ingress" {
		serviceType = data.Type
	} else {
		serviceType = "ClusterIP"
	}
	//组装ServiceCreate类型的数据
	sc := &ServiceCreate{
		Name:       getServiceName(data.Name),
		Namespace:  data.Namespace,
		Type:       serviceType,
		Port:       data.Port,
		TargetPort: data.ContainerPort,
		NodePort:   data.NodePort,
		Label:      data.Label,
	}
	// 创建service
	if err := Service.CreateService(client, sc); err != nil {
		return err
	}

	// 组装IngressCreate类型的数据，创建ingress，只有ingress类型的workflow才有ingress资源，所以 这里做了一层判断
	if data.Type == "Ingress" {
		ic := &IngressCreate{
			Name:      getIngressName(data.Name),
			Namespace: data.Namespace,
			Label:     data.Label,
			Hosts:     data.Hosts,
		}
		// 创建ingress
		if err := Ingress.CreateIngress(client, ic); err != nil {
			return err
		}
	}
	return nil
}

// 删除k8s资源 deployment service ingress
func deleteWorkFlowResource(client *kubernetes.Clientset, workflow *model.Workflow) (err error) {
	// 删除deployment
	if err := Deployment.DeleteDeployment(client, workflow.Name, workflow.Namespace); err != nil {
		return err
	}

	// 删除service
	if err := Service.DeleteService(client, getServiceName(workflow.Name), workflow.Namespace); err != nil {
		return err
	}

	// 删除ingress，这里多了一层判断，因为只有type为ingress的workflow才有ingress资源
	if workflow.Type == "Ingress" {
		if err := Ingress.DeleteIngress(client, getIngressName(workflow.Name), workflow.Namespace); err != nil {
			return err
		}
	}
	return nil
}

// workflow名字转换成service名字，添加-svc后缀
func getServiceName(workflowName string) (serviceName string) {
	return workflowName + "-svc"
}

// workflow名字转换成ingress名字，添加-ing后缀
func getIngressName(workflowName string) (ingressName string) {
	return workflowName + "-ing"
}