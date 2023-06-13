package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8sManagerApi/config"
)

var Pod = &pod{}

type pod struct {
}

// PodsResp 定义列表的返回内容，Items是Pod元素列表，Total是元素的数量
type PodsResp struct {
	Total int          `json:"total"`
	Items []corev1.Pod `json:"items"`
}

// PodsNp 获取每个namespace中pod数量，返回数据的结构体
type PodsNp struct {
	Namespace string `json:"namespace"`
	PodNum    int    `json:"pod_num"`
}

// GetPods 获取pod列表，支持过滤、排序、分页
func (p *pod) GetPods(client *kubernetes.Clientset, filterName, namespace string, limit, page int) (podsRest *PodsResp, err error) {

	// context.TODO()用于声明一个空的context上下文，用于List方法内设置这个请求的超时（源码），这里的常用用法
	// metav1.ListOptions{}用于过滤List数据，如使用label，field等
	// kubectl get services --all-namespaces --field-seletor metadata.namespace != default
	podList, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("获取Pod列表失败, %v\n", err.Error())
		return nil, errors.New("获取Pod列表失败," + err.Error())
	}
	// 实例化dataSelector结构体，组装数据
	selectableData := &dataSelector{
		GenericDataList: p.toCells(podList.Items),
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
	// 将DataCell类型转成Pod
	pods := p.fromCells(data.GenericDataList)

	/*
		// 数据处理后的数据和原始数据的比较
		// 处理后的数据
		fmt.Println("处理后的数据：")
		for _, pod := range pods {
			fmt.Println(pod.Name, pod.CreationTimestamp.Time)
		}
		// 原始数据
		fmt.Println("原始数据：")
		for _, pod := range podList.Items {
			fmt.Println(pod.Name, pod.CreationTimestamp.Time)
		}
	*/

	// 拼接返回数据
	podsRest = &PodsResp{
		Total: total,
		Items: pods,
	}
	return podsRest, nil
}

// GetPodDetail 获取Pod详情
func (p *pod) GetPodDetail(client *kubernetes.Clientset, podName, namespace string) (pod *corev1.Pod, err error) {
	pod, err = client.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("获取Pod详情失败, %v\n", err.Error())
		return nil, errors.New("获取Pod详情失败, " + err.Error())
	}
	return pod, nil
}

// DeletePod 删除Pod
func (p *pod) DeletePod(client *kubernetes.Clientset, podName, namespace string) (err error) {
	err = client.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil {
		fmt.Printf("删除Pod详情失败, %v\n", err.Error())
		return errors.New("删除Pod详情失败, " + err.Error())
	}
	return nil
}

// UpdatePod 更新Pod
func (p *pod) UpdatePod(client *kubernetes.Clientset, namespace, content string) (err error) {
	pod := &corev1.Pod{}
	err = json.Unmarshal([]byte(content), pod)
	if err != nil {
		fmt.Printf("反序列化失败, %v\n", err.Error())
		return errors.New("反序列化失败, " + err.Error())
	}
	_, err = client.CoreV1().Pods(namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		fmt.Printf("更新Pod失败, %v\n", err.Error())
		return errors.New("更新Pod失败, " + err.Error())
	}
	return nil
}

// GetPodContainer 获取Pod的容器名列表
func (p *pod) GetPodContainer(client *kubernetes.Clientset, podName, namespace string) (containers []string, err error) {
	pod, err := client.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("获取Pod详情失败, %v\n", err.Error())
		return nil, errors.New("获取Pod详情失败, " + err.Error())
	}
	for _, container := range pod.Spec.Containers {
		containers = append(containers, container.Name)
	}
	return containers, nil
}

// GetPodLog 获取容器的日志
func (p *pod) GetPodLog(client *kubernetes.Clientset, containerName, podName, namespace string) (Log string, err error) {
	// 设置日志的配置，容器名，获取的内容的配置
	lineLimit := int64(config.Conf.PodLogLine)
	option := &corev1.PodLogOptions{
		Container: containerName,
		TailLines: &lineLimit,
	}
	// 获取一个request实例
	req := client.CoreV1().Pods(namespace).GetLogs(podName, option)
	// 调用GetLogs方法 发起Stream连接，得到Response body
	logs, err := req.Stream(context.TODO())
	if err != nil {
		fmt.Printf("获取Pod日志失败, %v\n", err.Error())
		return "", errors.New("获取Pod日志失败, " + err.Error())
	}
	defer logs.Close()
	// 读取流中的数据，将response body 写入到缓冲区，目的是为了转换成string类型
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, logs)
	if err != nil {
		fmt.Printf("复制Pod日志失败, %v\n", err.Error())
		return "", errors.New("复制Pod日志失败, " + err.Error())
	}
	return buf.String(), nil
}

// GetPodNumPerNp 获取每个namespace的pod数量
func (p *pod) GetPodNumPerNp(client *kubernetes.Clientset) (podsNps []*PodsNp, err error) {
	// 获取Namespace列表
	namespaceList, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("获取Namespace列表失败, %v\n", err.Error())
		return nil, errors.New("获取Namespace列表失败, " + err.Error())
	}
	for _, namespace := range namespaceList.Items {
		// 获取pod列表
		podList, err := client.CoreV1().Pods(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Printf("获取Pod列表失败, %v\n", err.Error())
			return nil, errors.New("获取Pod列表失败, " + err.Error())
		}
		// 组装数据
		podsNp := &PodsNp{
			Namespace: namespace.Name,
			PodNum:    len(podList.Items),
		}
		// 添加数据到podsNps中
		podsNps = append(podsNps, podsNp)
	}
	return podsNps, nil
}

// 类型转换的方法，corev1.Pod -> DataCell, DataCell -> corev1.Pod
// toCells corev1.Pod -> DataCell
func (p *pod) toCells(pods []corev1.Pod) []DataCell {
	cells := make([]DataCell, len(pods))
	for i := range pods {
		cells[i] = podCell(pods[i])
	}
	return cells
}

// fromCells DataCell -> corev1.Pod
func (p *pod) fromCells(cells []DataCell) []corev1.Pod {
	pods := make([]corev1.Pod, len(cells))
	for i := range cells {
		//  cells[i].(podCell) 是将DataCell类型转成podCell
		pods[i] = corev1.Pod(cells[i].(podCell))
	}
	return pods
}