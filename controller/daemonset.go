package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"k8sManagerApi/service"
	"net/http"
)

var DaemonSet daemonSet

type daemonSet struct{}

// GetDaemonSetsHandler 获取DaemonSet列表，支持分页、过滤、排序
func (d *daemonSet) GetDaemonSetsHandler(ctx *gin.Context) {
	//	处理请求参数
	// 匿名结构体，用于定义传入参数，get请求为form格式，其它请求为json格式
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
		Cluster    string `form:"cluster"`
	})
	// form格式使用Bind方法，json格式使用SholdBindJson方法
	if err := ctx.Bind(params); err != nil {
		zap.L().Error(fmt.Sprintf("Bind绑定参数失败, %v", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Bind绑定参数失败" + err.Error(),
			"data": nil,
		})
		return
	}
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.DaemonSet.GetDaemonSets(client, params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success, 获取DaemonSet列表成功",
		"data": data,
	})
}

// GetDaemonSetDetailHandler 获取DaemonSet详情
func (d *daemonSet) GetDaemonSetDetailHandler(ctx *gin.Context) {
	//    处理请求参数
	// 匿名结构体，用于定义传入参数，get请求为form格式，其它请求为json格式
	params := new(struct {
		DaemonSetName string `form:"daemonSet_name"`
		Namespace     string `form:"namespace"`
		Cluster       string `form:"cluster"`
	})
	// form格式使用Bind方法，json格式使用SholdBindJson方法
	if err := ctx.Bind(params); err != nil {
		zap.L().Error(fmt.Sprintf("Bind绑定参数失败, %v", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Bind绑定参数失败" + err.Error(),
			"data": nil,
		})
		return
	}
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.DaemonSet.GetDaemonSetDetail(client, params.DaemonSetName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success, 获取DaemonSet详情成功",
		"data": data,
	})
}

// DeleteDaemonSetHandler 删除DaemonSet
func (d *daemonSet) DeleteDaemonSetHandler(ctx *gin.Context) {
	params := new(struct {
		DaemonSetName string `json:"daemonSet_name"`
		Namespace     string `json:"namespace"`
		Cluster       string `json:"cluster"`
	})
	// form格式使用Bind方法，json格式使用SholdBindJson方法
	if err := ctx.ShouldBindJSON(params); err != nil {
		zap.L().Error(fmt.Sprintf("Bind绑定参数失败, %v", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Bind绑定参数失败" + err.Error(),
			"data": nil,
		})
		return
	}
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	// 调用Service层方法进行删除
	if err := service.DaemonSet.DeleteDaemonSet(client, params.DaemonSetName, params.Namespace); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success, 删除DaemonSet成功",
		"data": nil,
	})
}

// UpdateDaemonSetHandler 更新DaemonSet
func (d *daemonSet) UpdateDaemonSetHandler(ctx *gin.Context) {
	params := new(struct {
		Namespace string `json:"namespace"`
		Content   string `json:"content"`
		Cluster   string `json:"cluster"`
	})
	// form格式使用Bind方法，json格式使用SholdBindJson方法
	if err := ctx.ShouldBindJSON(params); err != nil {
		zap.L().Error(fmt.Sprintf("Bind绑定参数失败, %v", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Bind绑定参数失败" + err.Error(),
			"data": nil,
		})
		return
	}
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err = service.DaemonSet.UpdateDaemonSet(client, params.Namespace, params.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success, 更新DaemonSet成功",
		"data": nil,
	})
}

// CreateDaemonSetHandler 创建DaemonSet
func (d *daemonSet) CreateDaemonSetHandler(ctx *gin.Context) {
	var (
		daemonSetCreate = new(service.DaemonSetCreate)
		err             error
	)
	if err = ctx.ShouldBindJSON(daemonSetCreate); err != nil {
		zap.L().Error(fmt.Sprintf("Bind绑定参数失败, %v", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Bind绑定参数失败" + err.Error(),
			"data": nil,
		})
		return
	}
	client, err := service.K8s.GetClient(daemonSetCreate.Cluster)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	// 调用Service层方法进行创建
	if err = service.DaemonSet.CreateDaemonSet(client, daemonSetCreate); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success, 创建Deployment成功",
		"data": nil,
	})
}

// GetDaemonSetNumPerNpHandler 获取每个Namespace的DaemonSet的数量
func (d *daemonSet) GetDaemonSetNumPerNpHandler(ctx *gin.Context) {
	params := new(struct {
		Cluster string `form:"cluster"`
	})
	if err := ctx.Bind(params); err != nil {
		zap.L().Error(fmt.Sprintf("获取参数失败, %v", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Bind绑定参数失败" + err.Error(),
			"data": nil,
		})
		return
	}
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.DaemonSet.GetDaemonSetNumPerNp(client)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success, 获取每个namespace的DaemonSet数量成功",
		"data": data,
	})
}