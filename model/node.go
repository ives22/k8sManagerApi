package model

import (
	"gorm.io/gorm"
	"time"
)

type Node struct {
	ID        uint           `json:"id" gorm:"primary_key"`
	CreatedAt *time.Time     `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Cluster        string `json:"cluster"`                      // 所属集群
	HostName       string `json:"name" gorm:"column:host_name"` // 主机名
	IP             string `json:"ip"`                           // 服务器IP地址
	Master         uint   `json:"master"`                       // 是否为master 1表示master 0表示work节点
	CPU            int    `json:"cpu"`                          // CPU核数
	Memory         string `json:"memory"`                       // 内存 3484264Ki
	System         string `json:"system"`                       // 服务器系统  linux
	OsImage        string `json:"os_image" gorm:"os_image"`     // 系统版本 CentOS Linux 7 (Core)
	Arch           string `json:"arch"`                         // 服务器架构 amd64
	KernelVersion  string `json:"kernel_version"`               // 服务器系统内核版本
	KubeletVersion string `json:"kubelet_version"`              // kubelet版本

}

func (*Node) TableName() string {
	return "node"
}

/*
CREATE TABLE `node` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `cluster` VARCHAR(255) DEFAULT NULL,
  `host_name` VARCHAR(255) DEFAULT NULL,
  `ip` VARCHAR(64) DEFAULT NULL,
  `master` TINYINT(1) DEFAULT NULL,
  `cpu` int(32) DEFAULT NULL,
  `memory` VARCHAR(255) DEFAULT NULL,
  `system` VARCHAR(255) DEFAULT NULL,
  `os_image` VARCHAR(255) DEFAULT NULL,
  `arch` VARCHAR(255) DEFAULT NULL,
  `kernel_version` VARCHAR(255) DEFAULT NULL,
  `kubelet_version` VARCHAR(255) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_node_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2291 CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
*/