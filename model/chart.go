package model

import (
	"gorm.io/gorm"
	"time"
)

type Chart struct {
	ID        uint           `json:"id" gorm:"primary_key"`
	CreatedAt *time.Time     `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name     string `json:"name"`
	FileName string `json:"file_name" gorm:"column:file_name"`
	IconUrl  string `json:"icon_url" gorm:"column:icon_url"`
	Version  string `json:"version" gorm:"column:version"`
	Describe string `json:"describe" gorm:"column:describe"`
}

func (*Chart) TableName() string {
	return "helm_chart"
}

/*
CREATE TABLE `helm_chart` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
	`name` VARCHAR(256) DEFAULT NULL,
  `file_name` VARCHAR(256) DEFAULT NULL,
  `icon_url` VARCHAR(256) DEFAULT NULL,
  `version` VARCHAR(256) DEFAULT NULL,
  `describe` VARCHAR(256) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_helm_chart_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2291 DEFAULT CHARSET=utf8mb4;
*/