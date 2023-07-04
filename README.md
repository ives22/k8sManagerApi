

## 建表语句
```sql
CREATE TABLE `helm_chart` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
	`name` VARCHAR(255) DEFAULT NULL,
  `file_name` VARCHAR(255) DEFAULT NULL,
  `icon_url` VARCHAR(255) DEFAULT NULL,
  `version` VARCHAR(255) DEFAULT NULL,
  `describe` VARCHAR(255) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_helm_chart_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2291 CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


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

```