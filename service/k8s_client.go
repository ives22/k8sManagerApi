package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8sManagerApi/config"
)

// 用于初始化k8是clientset

var K8s k8s

type k8s struct {
	ClientMap   map[string]*kubernetes.Clientset
	KubeConfMap map[string]string
}

// GetClient 获取client对象
func (k *k8s) GetClient(cluster string) (*kubernetes.Clientset, error) {
	client, ok := k.ClientMap[cluster]
	if !ok {
		fmt.Printf("集群不存在: %s，无法获取client\n", cluster)
		return nil, errors.New(fmt.Sprintf("集群不存在: %s, 无法获取client", cluster))
	}
	return client, nil
}

// Init 初始化
func (k *k8s) Init() {
	mp := map[string]string{}
	k.ClientMap = map[string]*kubernetes.Clientset{}
	if err := json.Unmarshal([]byte(config.KubeConfigs), &mp); err != nil {
		panic(fmt.Sprintf("Kubeconfigs反序列化失败 %v\n", err))
	}
	k.KubeConfMap = mp
	// 根据配置文件中配置的多个集群，循环进行初始化
	for key, value := range mp {
		conf, err := clientcmd.BuildConfigFromFlags("", value)
		if err != nil {
			panic(fmt.Sprintf("集群%s: 创建K8s 配置失败 %v", key, err))
		}
		clientSet, err := kubernetes.NewForConfig(conf)
		if err != nil {
			panic(fmt.Sprintf("集群%s: 创建K8s client失败 %v", key, err))
		}
		k.ClientMap[key] = clientSet
		fmt.Printf("集群%s: 创建K8s client成功\n", key)
	}
}