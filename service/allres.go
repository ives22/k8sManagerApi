package service

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sync"
)

var AllRes allRes

type allRes struct{}

// 定义一个全局互斥锁
var mt sync.Mutex

// GetAllNum 获取集群的所有资源
func (a *allRes) GetAllNum(client *kubernetes.Clientset) (map[string]int, []error) {
	var wg sync.WaitGroup
	wg.Add(14)
	errs := make([]error, 0)
	data := make(map[string]int, 0)

	// 获取所有的node节点
	go func() {
		list, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "Nodes", len(list.Items))
		wg.Done()
	}()
	// 获取所有的Namespace
	go func() {
		list, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "Namespaces", len(list.Items))
		wg.Done()
	}()
	// 获取所有的PV
	go func() {
		list, err := client.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "PVs", len(list.Items))
		wg.Done()
	}()
	// 获取所有的PVC
	go func() {
		list, err := client.CoreV1().PersistentVolumeClaims("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "PVCs", len(list.Items))
		wg.Done()
	}()
	// 获取所有的Service
	go func() {
		list, err := client.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "Services", len(list.Items))
		wg.Done()
	}()
	// 获取所有的Ingress
	go func() {
		list, err := client.NetworkingV1().Ingresses("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "Ingresses", len(list.Items))
		wg.Done()
	}()
	// 获取所有的Deployment
	go func() {
		list, err := client.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "Deployments", len(list.Items))
		wg.Done()
	}()
	// 获取所有的DaemonSet
	go func() {
		list, err := client.AppsV1().DaemonSets("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "DaemonSets", len(list.Items))
		wg.Done()
	}()
	// 获取所有的StatefulSet
	go func() {
		list, err := client.AppsV1().StatefulSets("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "StatefulSets", len(list.Items))
		wg.Done()
	}()
	// 获取所有的Job
	go func() {
		list, err := client.BatchV1().Jobs("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "Jobs", len(list.Items))
		wg.Done()
	}()
	// 获取所有的CronJobs
	go func() {
		list, err := client.BatchV1().CronJobs("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "CronJobs", len(list.Items))
		wg.Done()
	}()
	// 获取所有的Pod
	go func() {
		list, err := client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "Pods", len(list.Items))
		wg.Done()
	}()
	// 获取所有的Secrets
	go func() {
		list, err := client.CoreV1().Secrets("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "Secrets", len(list.Items))
		wg.Done()
	}()
	// 获取所有的ConfigMap
	go func() {
		list, err := client.CoreV1().ConfigMaps("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "ConfigMaps", len(list.Items))
		wg.Done()
	}()

	wg.Wait()
	return data, nil
}

func addMap(mp map[string]int, resource string, num int) {
	mt.Lock()
	defer mt.Unlock()
	mp[resource] = num
}

// 获取每个节点的所有pod信息