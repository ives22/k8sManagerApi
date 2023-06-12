package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	nwv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var Ingress ingress

type ingress struct{}

// IngressResp 定义列表的返回内容，Items是Ingress元素列表，Total是元素的数量
type IngressResp struct {
	Total int            `json:"total"`
	Items []nwv1.Ingress `json:"items"`
}

// IngressCreate 用于创建ingress需要的参数属性的定义
type IngressCreate struct {
	Name        string                 `json:"name"`
	Namespace   string                 `json:"namespace"`
	Annotations map[string]string      `json:"annotations"`
	Label       map[string]string      `json:"label"`
	Hosts       map[string][]*HttpPath `json:"hosts"`
	Cluster     string                 `json:"cluster"`
}

// HttpPath 定义ingress的path结构体
type HttpPath struct {
	Path        string        `json:"path"`
	PathType    nwv1.PathType `json:"path_type"`
	ServiceName string        `json:"service_name"`
	ServicePort int32         `json:"service_port"`
}

// IngressNp 用于返回namespace中service的数量
type IngressNp struct {
	Namespace  string `json:"namespace"`
	IngressNum int    `json:"ingress_num"`
}

// GetIngress 获取Ingress列表
func (i *ingress) GetIngress(client *kubernetes.Clientset, filterName, namespace string, limit, page int) (ingressRest *IngressResp, err error) {
	// 获取IngressList类型的Ingress列表
	serviceList, err := client.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("获取Ingress列表失败, %v\n", err.Error())
		return nil, errors.New("获取Ingress列表失败," + err.Error())
	}
	// 实例化dataSelector结构体，组装数据
	selectableData := &dataSelector{
		GenericDataList: i.toCells(serviceList.Items),
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
	//将[]DataCell类型的Ingress列表转为列表 nwv1.Ingress 类型
	ingressS := i.fromCells(data.GenericDataList)
	// 拼接返回数据
	ingressRest = &IngressResp{
		Total: total,
		Items: ingressS,
	}
	return ingressRest, nil
}

// GetIngressDetail 获取Ingress详情
func (i *ingress) GetIngressDetail(client *kubernetes.Clientset, ingressName, namespace string) (service *nwv1.Ingress, err error) {
	service, err = client.NetworkingV1().Ingresses(namespace).Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("获取Ingress详情失败, %v\n", err.Error())
		return nil, errors.New("获取Ingress详情失败, " + err.Error())
	}
	return service, nil
}

// CreateIngress 创建Ingress，接收IngressCreate对象
func (i *ingress) CreateIngress(client *kubernetes.Clientset, data *IngressCreate) (err error) {
	// 声明nwv1.IngressRule和nwv1.HTTPIngressPath变量，后面组装数据用
	var ingressRules []nwv1.IngressRule
	var httpIngressPaths []nwv1.HTTPIngressPath
	// 将data中的数据组装为nwv1.Ingress对象
	ingressCreat := &nwv1.Ingress{
		// ObjectMeta中定义资源名、名称空间以及标签
		ObjectMeta: metav1.ObjectMeta{
			Name:        data.Name,
			Namespace:   data.Namespace,
			Labels:      data.Label,
			Annotations: data.Annotations,
		},
		Status: nwv1.IngressStatus{},
	}
	// 第一层for循环是将host组装成nwv1.IngressRule类型的对象
	// 一个host对应一个ingressrule，每个ingressrule中包含一个host和多个path
	for key, value := range data.Hosts {
		ir := nwv1.IngressRule{
			Host: key,
			// 这里先将nwv1.HTTPIngressRuleValue类型中的Paths置为空，后面组装好数据再赋值
			IngressRuleValue: nwv1.IngressRuleValue{
				HTTP: &nwv1.HTTPIngressRuleValue{
					Paths: nil,
				},
			},
		}
		// 重新初始化 httpIngressPaths 切片 httpIngressPaths 是一个全局变量，它在每次循环迭代时都会保留之前的值。因此，当进行第二次循环时，httpIngressPaths 中仍然保留了第一次循环的结果。
		httpIngressPaths = []nwv1.HTTPIngressPath{}
		// 第二层for循环是将path组装成nwv1.HTTPIngressPath类型的对象
		for _, httpPath := range value {
			hip := nwv1.HTTPIngressPath{
				Path:     httpPath.Path,
				PathType: &httpPath.PathType,
				Backend: nwv1.IngressBackend{
					Service: &nwv1.IngressServiceBackend{
						Name: httpPath.ServiceName,
						Port: nwv1.ServiceBackendPort{
							Number: httpPath.ServicePort,
						},
					},
				},
			}
			// 将每个httpIngressPath对象组装成数据
			httpIngressPaths = append(httpIngressPaths, hip)
		}
		// 给Paths赋值，前面置为空了
		ir.IngressRuleValue.HTTP.Paths = httpIngressPaths
		//将每个ingressRule对象组装成数组，这个ingressRule对象就是IngressRule，每个元素是一个host和多个path
		ingressRules = append(ingressRules, ir)
	}
	//将ingressRules对象加入到ingress的规则中
	ingressCreat.Spec.Rules = ingressRules
	// 创建Ingress
	_, err = client.NetworkingV1().Ingresses(data.Namespace).Create(context.TODO(), ingressCreat, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("创建Ingress失败, %v\n", err.Error())
		return errors.New("创建Ingress失败, " + err.Error())
	}
	return nil
}

// DeleteIngress 删除Ingress
func (i *ingress) DeleteIngress(client *kubernetes.Clientset, ingressName, namespace string) (err error) {
	err = client.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), ingressName, metav1.DeleteOptions{})
	if err != nil {
		fmt.Printf("删除Ingress失败, %v\n", err.Error())
		return errors.New("删除Ingress失败, " + err.Error())
	}
	return nil
}

// UpdateIngress 更新Ingress
func (i *ingress) UpdateIngress(client *kubernetes.Clientset, namespace, content string) (err error) {
	var ingressS = &nwv1.Ingress{}
	err = json.Unmarshal([]byte(content), ingressS)
	if err != nil {
		fmt.Printf("反序列化失败 %v\n", err.Error())
		return errors.New("反序列化失败," + err.Error())
	}
	_, err = client.NetworkingV1().Ingresses(namespace).Update(context.TODO(), ingressS, metav1.UpdateOptions{})
	if err != nil {
		fmt.Printf("更新Ingress失败, %v\n", err.Error())
		return errors.New("更新Ingress失败, " + err.Error())
	}
	return nil
}

// 类型转换的方法，nwv1.Ingress-> DataCell, DataCell -> nwv1.Ingress
// toCells nwv1.Ingress-> DataCell
func (i *ingress) toCells(ingress []nwv1.Ingress) []DataCell {
	cells := make([]DataCell, len(ingress))
	for i := range ingress {
		cells[i] = ingressCell(ingress[i])
	}
	return cells
}

// fromCells DataCell -> nwv1.Ingress
func (i *ingress) fromCells(cells []DataCell) []nwv1.Ingress {
	ingressS := make([]nwv1.Ingress, len(cells))
	for i := range cells {
		//  cells[i].(daemonSetCell) 是将DataCell类型转成daemonSetCell
		ingressS[i] = nwv1.Ingress(cells[i].(ingressCell))
	}
	return ingressS
}