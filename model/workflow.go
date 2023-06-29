package model

import (
	"gorm.io/gorm"
	"time"
)

// Workflow 定义结构体，属性与mysql表字段对齐
type Workflow struct {
	//gorm:"primaryKey"用于声明主键
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt *time.Time     `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	Replicas   int32  `json:"replicas"`
	Deployment string `json:"deployment"`
	Service    string `json:"service"`
	Ingress    string `json:"ingress"`
	Type       string `json:"type" gorm:"column:type"`
	Cluster    string `json:"cluster"`
	//Type: clusterip nodeport ingress
}

// TableName 定义TableName方法，返回mysql表名，以此来定义mysql中的表名
func (*Workflow) TableName() string {
	return "workflow"
}

/*
CREATE TABLE `workflow` (
`id` int NOT NULL AUTO_INCREMENT,
`name` varchar(32) COLLATE utf8mb4_general_ci NOT NULL,
`namespace` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL,
`replicas` int DEFAULT NULL,
`deployment` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL,
`service` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL,
`ingress` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL,
`type` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL,
`cluster` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
`created_at` datetime DEFAULT NULL,
`updated_at` datetime DEFAULT NULL,
`deleted_at` datetime DEFAULT NULL,
PRIMARY KEY (`id`) USING BTREE,
UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
*/