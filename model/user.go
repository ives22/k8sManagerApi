package model

import (
	"gorm.io/gorm"
	"time"
)

// User 用户表
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt *time.Time     `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Username  string         `json:"username"`
	Password  string         `json:"password"`
}

// TableName 定义TableName方法，返回mysql表名，以此来定义mysql中的表名
func (*User) TableName() string {
	return "user"
}

/*
CREATE TABLE `user` (
`id` int NOT NULL AUTO_INCREMENT,
`username` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
`password` varchar(128) COLLATE utf8mb4_general_ci NOT NULL,
`created_at` datetime DEFAULT NULL,
`updated_at` datetime DEFAULT NULL,
`deleted_at` datetime DEFAULT NULL,
PRIMARY KEY (`id`) USING BTREE,
UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
*/