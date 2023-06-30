package service

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8sManagerApi/config"
	"k8sManagerApi/dao"
	"k8sManagerApi/model"
	"strconv"
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
		zap.L().Error("cluster not found", zap.String("cluster", cluster))
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

// AddClusterNodeInfo k8s client初始化完成后，获取该集群的node信息，并写入到数据库中
func (k *k8s) AddClusterNodeInfo() (err error) {
	// 循环每一个集群，获取node信息
	for cluster, client := range k.ClientMap {
		nodeList, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			zap.L().Warn(fmt.Sprintf("获取集群 %v 中的node列表失败, error: %v", cluster, err))
			return nil
		}
		for _, node := range nodeList.Items {
			var (
				hostname string
				ip       string
				isMaster uint
			)
			// 获取hostname 和 ip 地址
			for _, i := range node.Status.Addresses {
				if i.Type == "InternalIP" {
					ip = i.Address
				}
				if i.Type == "Hostname" {
					hostname = i.Address
				}
			}

			// 判断是否是master节点，1 是master节点 0 为work节点
			_, ok := node.Labels["node-role.kubernetes.io/master"]
			if !ok {
				isMaster = 0
			} else {
				isMaster = 1
			}

			// 整理数据
			cpu, _ := strconv.Atoi(node.Status.Capacity.Cpu().String())
			var nodeInfo = model.Node{
				Cluster:        cluster,
				HostName:       hostname,
				IP:             ip,
				Master:         isMaster,
				CPU:            cpu,
				Memory:         node.Status.Capacity.Memory().String(),
				System:         node.Status.NodeInfo.OperatingSystem,
				OsImage:        node.Status.NodeInfo.OSImage,
				Arch:           node.Status.NodeInfo.Architecture,
				KernelVersion:  node.Status.NodeInfo.KernelVersion,
				KubeletVersion: node.Status.NodeInfo.KubeletVersion,
			}

			// 调用数据库写入数据库中
			err := dao.Node.Add(&nodeInfo)
			if err != nil {
				zap.L().Error(fmt.Sprintf("Error creating node, %v", err))
			}
		}
	}
	zap.L().Info("Add cluster node info to database successfully.")
	return nil
}

// Init 初始化
func (k *k8s) Init() {
	k.ClientMap = map[string]*kubernetes.Clientset{}
	// 根据配置文件中的多个集群，循环进行初始化
	for _, cluster := range config.Conf.KubeConfigs {
		conf, err := clientcmd.BuildConfigFromFlags("", cluster.Path)
		if err != nil {
			zap.L().Error("create k8s client config failed", zap.String("cluster", cluster.Name))
			panic(fmt.Sprintf("集群%s: 创建K8s 配置失败 %v", cluster.Name, cluster.Path))
		}
		clientSet, err := kubernetes.NewForConfig(conf)
		if err != nil {
			zap.L().Error("create k8s client failed", zap.String("cluster", cluster.Name))
			panic(fmt.Sprintf("集群%s: 创建K8s client失败 %v", cluster.Name, cluster.Path))
		}
		k.ClientMap[cluster.Name] = clientSet
		zap.L().Info("create k8s client successfully", zap.String("cluster", cluster.Name))
	}

	_ = k.AddClusterNodeInfo()
}