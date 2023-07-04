package model

import (
	"gorm.io/gorm"
	"time"
)

type PodInfo struct {
	ID        uint           `json:"id" gorm:"primary_key"`
	CreatedAt *time.Time     `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Cluster      string    `json:"cluster"`                             // 所属集群
	PodName      string    `json:"pod_name" gorm:"column:pod_name"`     // pod名字
	HostIP       string    `json:"host_ip" gorm:"column:host_ip"`       // 所在节点IP
	PodIP        string    `json:"pod_ip" gorm:"pod_ip"`                // Pod IP
	Status       string    `json:"status"`                              // pod状态
	CreationTime time.Time `json:"creation_time" gorm:"type:timestamp"` // Pod的创建时间
}

func (*PodInfo) TableName() string {
	return "pod_info"
}

/*
CREATE TABLE `pod_info` (
  `id` int NOT NULL AUTO_INCREMENT,
  `cluster` VARCHAR(255) DEFAULT NULL,
  `pod_name` VARCHAR(255) DEFAULT NULL,
  `host_ip` VARCHAR(64) DEFAULT NULL,
  `pod_ip` VARCHAR(64) DEFAULT NULL,
  `status` VARCHAR(64) DEFAULT NULL,
  `creation_time` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_pod_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2291 CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
*/