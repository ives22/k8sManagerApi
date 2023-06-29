package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Router 初始化router类型对象，首字母大写，用于跨包调用
var Router router

// 声明一个router的结构体
type router struct{}

// InitApiRouter 初始化API路由
func (r *router) InitApiRouter(router *gin.Engine) {
	// GET请求，路径为"/testapi"，处理函数返回JSON数据
	router.GET("/testapi", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "test message",
			"data": "ok",
		})
	})

	// 以下为Kubernetes相关的路由和处理函数

	// 获取所有Pods的路由，GET请求，路径为"/api/k8s/pods"，处理函数为Pod.GetPodsHandler
	router.GET("/api/k8s/pods", Pod.GetPodsHandler)
	// 获取单个Pod详情的路由，GET请求，路径为"/api/k8s/pod/detail"，处理函数为Pod.GetPodDetailHandler
	router.GET("/api/k8s/pod/detail", Pod.GetPodDetailHandler)
	// 获取Pod容器的路由，GET请求，路径为"/api/k8s/pod/container"，处理函数为Pod.GetPodContainerHandler
	router.GET("/api/k8s/pod/container", Pod.GetPodContainerHandler)
	// 获取Pod日志的路由，GET请求，路径为"/api/k8s/pod/log"，处理函数为Pod.GetPodLogHandler
	router.GET("/api/k8s/pod/log", Pod.GetPodLogHandler)
	// 获取Pod数量的路由，GET请求，路径为"/api/k8s/pod/numnp"，处理函数为Pod.GetPodNumberHandler
	router.GET("/api/k8s/pod/numnp", Pod.GetPodNumberHandler)
	// 更新Pod的路由，PUT请求，路径为"/api/k8s/pod/update"，处理函数为Pod.UpdatePodHandler
	router.PUT("/api/k8s/pod/update", Pod.UpdatePodHandler)
	// 删除Pod的路由，DELETE请求，路径为"/api/k8s/pod/del"，处理函数为Pod.DeletePodHandler
	router.DELETE("/api/k8s/pod/del", Pod.DeletePodHandler)

	// 以下为Deployment相关的路由和处理函数
	// 获取所有Deployments的路由，GET请求，路径为"/api/k8s/deployments"，处理函数为Deployment.GetDeploymentsHandler
	router.GET("/api/k8s/deployments", Deployment.GetDeploymentsHandler)
	// 获取单个Deployment详情的路由，GET请求，路径为"/api/k8s/deployment/detail"，处理函数为Deployment.GetDeploymentDetailHandler
	router.GET("/api/k8s/deployment/detail", Deployment.GetDeploymentDetailHandler)
	// 获取Deployment数量的路由，GET请求，路径为"/api/k8s/deployment/numnp"，处理函数为Deployment.GetDeployNumPerNpHandler
	router.GET("/api/k8s/deployment/numnp", Deployment.GetDeployNumPerNpHandler)
	// 更新Deployment副本数量的路由，PUT请求，路径为"/api/k8s/deployment/scale"，处理函数为Deployment.SetDeploymentReplicasHandler
	router.PUT("/api/k8s/deployment/scale", Deployment.SetDeploymentReplicasHandler)
	// 重启Deployment的路由，PUT请求，路径为"/api/k8s/deployment/restart", 处理函数为Deployment.RestartDeploymentHandler
	router.PUT("/api/k8s/deployment/restart", Deployment.RestartDeploymentHandler)
	// 更新Deployment的路由，PUT请求，路径为"/api/k8s/deployment/update"，处理函数为Deployment.UpdateDeploymentHandler
	router.PUT("/api/k8s/deployment/update", Deployment.UpdateDeploymentHandler)
	// 创建Deployment的路由，POST请求，路径为"/api/k8s/deployment/create", 处理函数为Deployment.CreateDeploymentHandler
	router.POST("/api/k8s/deployment/create", Deployment.CreateDeploymentHandler)
	// 删除Deployment的路由，DELETE请求，路径为"/api/k8s/deployment/del"，处理函数为Deployment.DeleteDeploymentHandler
	router.DELETE("/api/k8s/deployment/del", Deployment.DeleteDeploymentHandler)

	// 以下为DaemonSet相关的路由和处理函数
	router.GET("/api/k8s/daemonSets", DaemonSet.GetDaemonSetsHandler)
	router.GET("/api/k8s/daemonSet/detail", DaemonSet.GetDaemonSetDetailHandler)
	router.GET("/api/k8s/daemonSet/numnp", DaemonSet.GetDaemonSetNumPerNpHandler)
	router.DELETE("/api/k8s/daemonSet/del", DaemonSet.DeleteDaemonSetHandler)
	router.PUT("/api/k8s/daemonSet/update", DaemonSet.UpdateDaemonSetHandler)
	router.POST("/api/k8s/daemonSet/create", DaemonSet.CreateDaemonSetHandler)

	// 以下为StatefulSet相关的路由和处理函数
	router.GET("/api/k8s/statefulSets", StatefulSet.GetStatefulSetsHandler)
	router.GET("/api/k8s/statefulSet/detail", StatefulSet.GetStatefulSetDetailHandler)
	router.GET("/api/k8s/statefulSet/numnp", StatefulSet.GetStatefulSetNumPerNpHandler)
	router.DELETE("/api/k8s/statefulSet/del", StatefulSet.DeleteStatefulSetHandler)
	router.PUT("/api/k8s/statefulSet/update", StatefulSet.UpdateStatefulSetHandler)

	// 以下是Node相关的路由和处理函数
	router.GET("/api/k8s/nodes", Node.GetNodesHandler)
	router.GET("/api/k8s/node/detail", Node.GetNodeDetailHandler)

	// 以下是Namespace相关的路由和处理函数
	router.GET("/api/k8s/namespaces", Namespace.GetNamespacesHandler)
	router.GET("/api/k8s/namespace/detail", Namespace.GetNamespaceDetailHandler)
	router.DELETE("/api/k8s/namespace/del", Namespace.DeleteNamespaceHandler)
	router.POST("/api/k8s/namespace/create", Namespace.CreateNamespaceHandler)

	// 以下是pv相关的路由和处理函数
	router.GET("/api/k8s/pvs", Pv.GetPvsHandler)
	router.GET("/api/k8s/pv/detail", Pv.GetPvDetailHandler)
	router.DELETE("/api/k8s/pv/del", Pv.DeletePvHandler)

	// 以下为Service相关的路由和处理函数
	router.GET("/api/k8s/services", Service.GetServicesHandler)
	router.GET("/api/k8s/service/detail", Service.GetServicesDetailHandler)
	//router.GET("/api/k8s/service/numnp", StatefulSet.GetStatefulSetNumPerNpHandler)
	router.POST("/api/k8s/service/create", Service.CreateServiceHandler)
	router.DELETE("/api/k8s/service/del", Service.DeleteServiceHandler)
	router.PUT("/api/k8s/service/update", Service.UpdateServiceHandler)

	// 以下为Ingress相关的路由和处理函数
	router.GET("/api/k8s/ingresses", Ingress.GetIngressHandler)
	router.GET("/api/k8s/ingress/detail", Ingress.GetIngressDetailHandler)
	router.POST("/api/k8s/ingress/create", Ingress.CreateIngressHandler)
	router.DELETE("/api/k8s/ingress/del", Ingress.DeleteIngressHandler)
	router.PUT("/api/k8s/ingress/update", Ingress.UpdateIngressHandler)

	// 以下为workflow相关的路由和处理函数
	router.GET("/api/k8s/workflows", Workflow.GetWorkflowsHandler)
	router.GET("/api/k8s/workflow/detail", Workflow.GetWorkflowDetailHandler)
	router.POST("/api/k8s/workflow/create", Workflow.CreateWorkflowHandler)
	router.DELETE("/api/k8s/workflow/del", Workflow.DeleteWorkflowHandler)

	// 以下是ConfigMap相关的路由和处理函数
	router.GET("/api/k8s/configmaps", ConfigMap.GetConfigMapsHandler)
	router.GET("/api/k8s/configmap/detail", ConfigMap.GetConfigMapDetailHandler)
	router.PUT("/api/k8s/configmap/update", ConfigMap.UpdateConfigMapHandler)
	router.DELETE("/api/k8s/configmap/del", ConfigMap.DeleteConfigMapHandler)
	// 以下是Secret相关的路由和处理函数
	router.GET("/api/k8s/secrets", Secret.GetSecretsHandler)
	router.GET("/api/k8s/secret/detail", Secret.GetSecretDetailHandler)
	router.PUT("/api/k8s/secret/update", Secret.UpdateSecretHandler)
	router.DELETE("/api/k8s/secret/del", Secret.DeleteSecretHandler)
	// 以下是PVC相关的路由和处理函数
	router.GET("/api/k8s/pvcs", Pvc.GetPvcsHandler)
	router.GET("/api/k8s/pvc/detail", Pvc.GetPvcDetailHandler)
	router.PUT("/api/k8s/pvc/update", Pvc.UpdatePvcHandler)
	router.DELETE("/api/k8s/pvc/del", Pvc.DeletePvcHandler)

	// 以下是登录注册相关的路由和处理函数
	router.POST("/api/login", Auth.LoginHandler)
	// 注册路由
	router.POST("/api/register", Auth.RegisterHandler)

	// 获取集群列表
	router.GET("/api/k8s/clusters", Cluster.GetClustersHandler)

	// 获取集群所有资源
	router.GET("/api/k8s/allres", AllRes.GetAllNumHandler)

	// 获取集群事件
	router.GET("/api/k8s/events", Event.GetEventsHandler)

	// helm 应用商店
	router.GET("/api/helmstore/releases", HelmStore.ListReleasesHandler)
	router.GET("/api/helmstore/release/detail", HelmStore.DetailReleaseHandler)
	router.POST("/api/helmstore/release/install", HelmStore.InstallReleaseHandler)
	router.DELETE("/api/helmstore/release/uninstall", HelmStore.UninstallReleaseHandler)
	router.GET("/api/helmstore/charts", HelmStore.GetChartsHandler)
	router.POST("/api/helmstore/chart/add", HelmStore.AddChartHandler)
	router.PUT("/api/helmstore/chart/update", HelmStore.UpdateChartHandler)
	router.DELETE("/api/helmstore/chart/del", HelmStore.DeleteChartHandler)
	router.POST("/api/helmstore/chartfile/upload", HelmStore.UploadChartFileHandler)
	router.DELETE("/api/helmstore/chartfile/del", HelmStore.DeleteChartFileHandler)

	// 用于服务器下载chart文件进行安装，下载文件的路由
	router.GET("/download/:filename", Download.DownloadHandler)

	// 以下为StatefulSet相关的路由和处理函数
	// 获取所有StatefulSets的路由，GET请求，路径为"/api/k8s/statefulsets"，处理函数为StatefulSet.GetStatefulSetsHandler
	//router.GET("/api/k8s/statefulsets", StatefulSet.GetStatefulSetsHandler)
	// 获取单个StatefulSet详情的路由，GET请求，路径为"/api/k8s/statefulset/detail"，处理函数为StatefulSet.GetStatefulSetDetailHandler
	//router.GET("/api/k8s/statefulset/detail", StatefulSet.GetStatefulSetDetailHandler)
}