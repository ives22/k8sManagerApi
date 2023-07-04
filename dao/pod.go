package dao

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"k8sManagerApi/db/mysql"
	"k8sManagerApi/model"
)

var PodInfo podInfo

type podInfo struct{}

// Add Pod信息
func (p *podInfo) Add(pod *model.PodInfo) (err error) {
	// 先判断Pod是否存在，如果存在则不添加
	existingPod := &model.PodInfo{}
	result := mysql.DB.Where("pod_name = ? and cluster = ? and pod_ip = ? and host_ip = ? and creation_time = ?",
		pod.PodName, pod.Cluster, pod.PodIP, pod.HostIP, pod.CreationTime).First(&existingPod)

	if result == nil {
		// 如果Pod已存在，便在判断Pod的状态是否相等，如果相等，则不创建，否则只是更新状态。
		res := mysql.DB.Where("pod_name = ? and cluster = ? and pod_ip = ? and host_ip = ? and creation_time = ? and status = ?",
			pod.PodName, pod.Cluster, pod.PodIP, pod.HostIP, pod.CreationTime, pod.Status).First(&existingPod)
		if res == nil {
			zap.L().Warn(fmt.Sprintf("pod已存在, cluster: %v, name: %v", pod.Cluster, pod.PodName))
			return nil
		}
		// 如果pod只是状态发生了改变，则只是修改状态即可。
		mysql.DB.Model(&existingPod).Update("status", pod.Status)

	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 查询出错
		zap.L().Error(fmt.Sprintf("查询Pod失败: %v", result.Error))
		return errors.New(fmt.Sprintf("查询Pod失败: %v", result.Error))
	}

	// 插入新节点
	tx := mysql.DB.Create(&pod)
	if tx.Error != nil {
		//zap.L().Error(fmt.Sprintf("新增Pod失败, %v", tx.Error))
		return errors.New(fmt.Sprintf("新增Pod失败, %v", tx.Error))
	}
	return nil
}

// Del 删除Pod信息
func (p *podInfo) Del(pod *model.PodInfo) (err error) {
	delPod := model.PodInfo{}
	mysql.DB.Debug().Where("pod_name = ? and cluster = ?", pod.PodName, pod.Cluster).Delete(&delPod)
	return nil
}