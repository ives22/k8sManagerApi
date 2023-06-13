package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8sManagerApi/config"
	"k8sManagerApi/controller"
	"k8sManagerApi/db/mysql"
	"k8sManagerApi/middle"
	"k8sManagerApi/service"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	// 初始化配置文件
	config.Init()

	//// 测试配置文件是否成功加载
	//fmt.Println("Listen Addr: ", config.Conf.ListenAddr)

	// 初始化数据库
	mysql.Init()

	// 初始化k8s client
	service.K8s.Init()
	//service.K8sA.Init()

	//	初始化gin
	r := gin.Default()

	// 注册中间件, 跨域配置
	r.Use(middle.Cors())
	// 注册中间件，加载jwt中间件
	r.Use(middle.JWTAuth())

	//	跨包调用router的初始化方法
	controller.Router.InitApiRouter(r)

	// 终端websocket
	go func() {
		http.HandleFunc(config.Conf.WSPath, service.Terminal.WebsocketHandler)
		http.ListenAndServe(config.Conf.WSAddr, nil)
	}()

	// event任务，用于监听event并写入数据库，这里传入的参数是集群名，与config配置文件中集群名对齐
	//go func() {
	//	service.Event.WatchEventTask("TST-1")
	//}()
	// event任务，用于监听event并写入数据库，循环从配置文件中读取集群名，启动g oroutin 任务
	for _, cInfo := range config.Conf.KubeConfigs {
		go func(clusterName string) {
			service.Event.WatchEventTask(clusterName)
		}(cInfo.Name)
	}

	// 数据库测试
	//data, _ := dao.User.GetUserByName("zhangsan")
	//fmt.Println("data: ", data)

	//	启动gin server
	//r.Run(config.ListerAddr)
	srv := &http.Server{
		Addr:    config.Conf.ListenAddr,
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("listen: %v\n", err)
		}
	}()
	// 等待中断信号，优雅关闭所有server
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	// 设置ctx超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// cancel用于释放ctx
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Gin Server关闭异常: %v\n", err)
	}
	fmt.Println("Gin Server关闭成功")
}