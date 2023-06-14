package service

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8sManagerApi/dao"
	"k8sManagerApi/model"
	"time"
)

var Event event

type event struct{}

// GetEvents 获取events列表
func (e *event) GetEvents(name, cluster string, page, limit int) (events *dao.Events, err error) {
	data, err := dao.Event.GetEvents(name, cluster, page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// WatchEventTask informer监听event
func (e *event) WatchEventTask(cluster string) {
	// 实例化 informerFactory
	informerFactory := informers.NewSharedInformerFactory(K8s.ClientMap[cluster], time.Minute)
	// 监听资源
	informer := informerFactory.Core().V1().Events()
	// 添加事件handler
	informer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				onAdd(obj, cluster)
			},
		},
	)
	// 处理启动和优雅关闭
	stopCh := make(chan struct{})
	defer close(stopCh)
	informerFactory.Start(stopCh)
	if !cache.WaitForCacheSync(stopCh, informer.Informer().HasSynced) {
		fmt.Println("同步cache超时")
		return
	}
	<-stopCh
	return
}

// onAdd 新增时落库
func onAdd(obj interface{}, cluster string) {
	// 断言
	event := obj.(*corev1.Event)
	// 判断是否重复
	_, has, err := dao.Event.HasEvent(
		event.InvolvedObject.Name,
		event.InvolvedObject.Kind,
		event.InvolvedObject.Namespace,
		event.Reason,
		event.CreationTimestamp.Time,
		cluster,
	)
	if err != nil {
		return
	}
	if has {
		//fmt.Printf("Event数据已存在, %s %s %s %s %v %s\n",
		//	event.InvolvedObject.Name,
		//	event.InvolvedObject.Kind,
		//	event.InvolvedObject.Namespace,
		//	event.Reason,
		//	event.CreationTimestamp.Time,
		//	cluster)
		return
	}
	// 组装数据
	data := &model.Event{
		Name:      event.InvolvedObject.Name,
		Kind:      event.InvolvedObject.Kind,
		Namespace: event.InvolvedObject.Namespace,
		Rtype:     event.Type,
		Reason:    event.Reason,
		Message:   event.Message,
		EventTime: &event.CreationTimestamp.Time,
		Cluster:   cluster,
	}
	// 数据库添加
	if err := dao.Event.Add(data); err != nil {
		return
	}
}