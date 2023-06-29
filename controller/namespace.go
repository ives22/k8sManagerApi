package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"k8sManagerApi/service"
	"net/http"
)

var Namespace namespace

type namespace struct{}

// GetNamespacesHandler 获取Namespace列表
func (n *namespace) GetNamespacesHandler(ctx *gin.Context) {
	//	处理请求参数
	// 匿名结构体，用于定义传入参数，get请求为form格式，其它请求为json格式
	params := new(struct {
		FilterName string `form:"filter_name"`
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
	data, err := service.Namespace.GetNamespaces(client, params.FilterName, params.Limit, params.Page)
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
		"msg":  "success, 获取Namespace列表成功",
		"data": data,
	})
}

// GetNamespaceDetailHandler 获取namespace详情
func (n *namespace) GetNamespaceDetailHandler(ctx *gin.Context) {
	// 匿名结构体，用于定义传入参数，get请求为form格式，其它请求为json格式
	params := new(struct {
		NamespaceName string `form:"namespace_name"`
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
	data, err := service.Namespace.GetNamespaceDetail(client, params.NamespaceName)
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
		"msg":  "success, 获取Namespace详情成功",
		"data": data,
	})
}

// CreateNamespaceHandler 创建namespace
func (n *namespace) CreateNamespaceHandler(ctx *gin.Context) {
	var (
		namespaceCreate = new(service.NamespaceCreate)
		err             error
	)
	// form格式使用Bind方法，json格式使用SholdBindJson方法
	if err = ctx.ShouldBindJSON(namespaceCreate); err != nil {
		zap.L().Error(fmt.Sprintf("Bind绑定参数失败, %v", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Bind绑定参数失败" + err.Error(),
			"data": nil,
		})
		return
	}
	client, err := service.K8s.GetClient(namespaceCreate.Cluster)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	// 调用Service层方法进行创建
	if err = service.Namespace.CreateNamespace(client, namespaceCreate); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success, 创建Namespace成功",
		"data": nil,
	})
}

// DeleteNamespaceHandler 删除namespace
func (n *namespace) DeleteNamespaceHandler(ctx *gin.Context) {
	params := new(struct {
		NamespaceName string `json:"namespace_name"`
		Cluster       string `json:"cluster"`
	})
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
	err = service.Namespace.DeleteNamespace(client, params.NamespaceName)
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
		"msg":  "success, 删除Namespace成功",
	})
}