package service

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"helm.sh/helm/v3/pkg/action"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"log"
	"os"
)

var HelmConfig helmConfig

type helmConfig struct {
	// 这种方式初始化行不通，比如有10个命名空间。
	// ActionConfigMap map[string]*action.Configuration
}

// GetAc 获取Helm action 配置
func (c *helmConfig) GetAc(cluster, namespace string) (*action.Configuration, error) {
	kubeconfig := K8s.GetClusterConf(cluster)
	fmt.Println("kube配置文件：", kubeconfig)
	if kubeconfig == "" {
		zap.L().Warn(fmt.Sprintf("集群不存在: %s, 无法获取client", cluster))
		return nil, errors.New(fmt.Sprintf("集群不存在: %s, 无法获取client", cluster))
	}

	// new一个actionConfig对象
	actionConfig := new(action.Configuration)
	cf := &genericclioptions.ConfigFlags{
		KubeConfig: &kubeconfig,
		Namespace:  &namespace,
	}
	if err := actionConfig.Init(cf, namespace, os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		zap.L().Error(fmt.Sprintf("actionConfig初始化失败, %v", err.Error()))
		return nil, errors.New("actionConfig初始化失败, " + err.Error())
	}
	return actionConfig, nil
}