package service

import (
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

// GetClusterConf 获取指定集群的配置文件
func (k *k8s) GetClusterConf(cluster string) (clusterConf string) {
	for _, conf := range config.Conf.KubeConfigs {
		if conf.Name == cluster {
			clusterConf = conf.Path
			return clusterConf
		}
	}
	return ""
}

// Init 初始化
func (k *k8s) Init() {
	k.ClientMap = map[string]*kubernetes.Clientset{}
	// 根据配置文件中的多个集群，循环进行初始化
	for _, cluster := range config.Conf.KubeConfigs {
		conf, err := clientcmd.BuildConfigFromFlags("", cluster.Path)
		if err != nil {
			panic(fmt.Sprintf("集群%s: 创建K8s 配置失败 %v", cluster.Name, cluster.Path))
		}
		clientSet, err := kubernetes.NewForConfig(conf)
		if err != nil {
			panic(fmt.Sprintf("集群%s: 创建K8s client失败 %v", cluster.Name, cluster.Path))
		}
		k.ClientMap[cluster.Name] = clientSet
		fmt.Printf("集群%s: 创建K8s client成功\n", cluster.Name)
	}
}