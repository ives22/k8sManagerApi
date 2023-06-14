package config

type ServerConfig struct {
	ListenAddr  string        `mapstructure:"listenAddr"`
	WSAddr      string        `mapstructure:"WSAddr"`
	WSPath      string        `mapstructure:"WSPath"`
	PodLogLine  int64         `mapstructure:"podLogTailLine"`
	UploadPath  string        `mapstructure:"uploadPath"`
	KubeConfigs []*Kubeconfig `mapstructure:"KubeConfigs"`
	MysqlInfo   *MysqlConfig  `mapstructure:"mysql"`
	LogConfig   *LogConfig    `mapstructure:"log"`
}

type Kubeconfig struct {
	Name string `mapstructure:"name"`
	Path string `mapstructure:"path"`
}

type MysqlConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Charset  string `mapstructure:"charset"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"fileName"`
	MaxSize    int    `mapstructure:"maxSize"`
	MaxAge     int    `mapstructure:"maxAge"`
	MaxBackups int    `mapstructure:"maxBackups"`
	Compress   bool   `mapstructure:"compress"`
}