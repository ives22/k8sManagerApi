package config

import "time"

const (
	ListerAddr = "0.0.0.0:9091"
	KubeConfig = "/Users/liyanjie/Documents/config"
	//KubeConfig = "/root/.kube/config"
	//验证多集群
	KubeConfigs = `{"TST-1":"/Users/liyanjie/Documents/config", "TST-2":"/Users/liyanjie/Documents/configw", "TST-3":"/Users/liyanjie/Documents/configyb"}`
	//KubeConfigs = `{"TST-1":"/root/.kube/config", "TST-2":"/root/.kube/config"}`
	// tail的日志行数
	// tail -n 2000
	PodLogTailLine = 2000

	// 数据库配置
	DBType     = "mysql"
	DBHost     = "124.71.33.240"
	DBPort     = 3306
	DBUser     = "root"
	DBPassword = "admin123"
	DBName     = "k8s"
	DBCharset  = "utf8"

	// 打印MySQL debug的sql日志
	DebugSQL = false

	// 连接池的配置
	MaxIdleConns = 10               // 最大空闲连接
	MaxOpenConns = 100              // 最大打开连接
	MaxLifeTime  = 30 * time.Second // 最大生存时间

	// websocket 配置
	WSAddr = "0.0.0.0:9092"
	WSPath = "/ws"

	// helm文件上传路径
	UploadPath = "/Users/liyanjie/Documents/"
	//UploadPath = "/root/charts/"
)