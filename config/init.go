package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Conf 全局变量，存放配置信息
var Conf = new(ServerConfig)

// GetEnvInfo 获取环境变量
func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func Init() {
	// 获取环境变量，根据环境变量去加载不同的配置文件, 如果环境变量K8SM_DEBUT的值为true，则加载本地配置文件，如果为false则加载生产环境的配置文件
	debug := GetEnvInfo("K8SM_DEBUT")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("etc/%s_pro.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("etc/%s_local.yaml", configFilePrefix)
	}

	v := viper.New()
	v.SetConfigFile(configFileName)

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			panic("Config file not found")
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("Fatal error config file, %s \n", err))
		}
	}

	// 将读取的配置信息保存至全局变量Conf
	if err := v.Unmarshal(Conf); err != nil {
		panic(fmt.Errorf("unmarshal conf failed, %v\n", err))
	}

	// 监听配置文件是否发生改变
	v.WatchConfig()
	// 如果发生了改变，需要同步给全局变量Conf
	v.OnConfigChange(func(in fsnotify.Event) {
		if err := v.Unmarshal(Conf); err != nil {
			panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
		}
	})
}