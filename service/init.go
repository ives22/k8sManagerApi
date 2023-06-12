package service

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8sManagerApi/config"
)

// 用于初始化k8是clientset

var K8sO k8so

type k8so struct {
	Clientset *kubernetes.Clientset
}

// Init 初始化
func (k *k8so) Init() {
	//  将kubeconfig格式化为rest.config类型的对象
	conf, err := clientcmd.BuildConfigFromFlags("", config.KubeConfig)
	if err != nil {
		panic("获取k8s client配置失败, " + err.Error())
	}
	//	通过config创建clientset
	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		panic("创建k8s client失败, " + err.Error())
	} else {
		fmt.Println("k8s client 初始化成功!")
	}
	k.Clientset = clientset
}