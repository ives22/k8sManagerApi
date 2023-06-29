package model

import (
	"gorm.io/gorm"
	"time"
)

type Event struct {
	ID        uint           `json:"id" gorm:"primary_key"`
	CreatedAt *time.Time     `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name      string     `json:"name"`
	Kind      string     `json:"kind"`
	Namespace string     `json:"namespace"`
	Rtype     string     `json:"rtype"`
	Reason    string     `json:"reason"`  // 事件原因
	Message   string     `json:"message"` // 事件描述
	EventTime *time.Time `json:"event_time"`
	Cluster   string     `json:"cluster"`
}

func (*Event) TableName() string {
	return "k8s_event"
}

/*
CREATE TABLE `k8s_event` (
`id` int(11) NOT NULL AUTO_INCREMENT,
`name` varchar(255) DEFAULT NULL,
`kind` varchar(255) DEFAULT NULL,
`namespace` varchar(255) DEFAULT NULL,
`rtype` varchar(255) DEFAULT NULL,
`reason` varchar(255) DEFAULT NULL,
`message` varchar(255) DEFAULT NULL,
`event_time` datetime DEFAULT NULL,
`cluster` varchar(64) DEFAULT NULL,
`created_at` datetime DEFAULT NULL,
`updated_at` datetime DEFAULT NULL,
`deleted_at` datetime DEFAULT NULL,
PRIMARY KEY (`id`),
KEY `idx_k8s_event_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2291 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
*/