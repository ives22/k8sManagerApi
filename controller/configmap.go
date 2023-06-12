package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8sManagerApi/service"
	"net/http"
)

var ConfigMap configMap

type configMap struct{}

// GetConfigMapsHandler 获取configmap列表，支持分页、过滤、排序
func (c *configMap) GetConfigMapsHandler(ctx *gin.Context) {
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
		fmt.Printf("绑定参数失败, %v\n", err.Error())
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
	data, err := service.ConfigMap.GetConfigMaps(client, params.Namespace, params.FilterName, params.Limit, params.Page)
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
		"msg":  "success, 获取ConfigMap列表成功",
		"data": data,
	})
}

// GetConfigMapDetailHandler 获取configMap详情
func (c *configMap) GetConfigMapDetailHandler(ctx *gin.Context) {
	params := new(struct {
		Namespace     string `form:"namespace"`
		ConfigmapName string `form:"configmap_name"`
		Cluster       string `form:"cluster"`
	})
	if err := ctx.Bind(params); err != nil {
		fmt.Printf("绑定参数失败, %v\n", err.Error())
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
	data, err := service.ConfigMap.GetConfigMapDetail(client, params.Namespace, params.ConfigmapName)
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
		"msg":  "success, 获取ConfigMap详情成功",
		"data": data,
	})
}

// UpdateConfigMapHandler 更新ConfigMap
func (c *configMap) UpdateConfigMapHandler(ctx *gin.Context) {
	params := new(struct {
		Namespace string `json:"namespace"`
		Content   string `json:"content"`
		Cluster   string `json:"cluster"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		fmt.Printf("绑定参数失败, %v\n", err.Error())
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
	err = service.ConfigMap.UpdateConfigMap(client, params.Namespace, params.Content)
	if err != nil {
		fmt.Printf("更新ConfigMap失败%v\n", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success, 更新ConfigMap成功",
		"data": nil,
	})
}

// DeleteConfigMapHandler 删除ConfigMap
func (c *configMap) DeleteConfigMapHandler(ctx *gin.Context) {
	params := new(struct {
		Namespace     string `json:"namespace"`
		ConfigmapName string `json:"configmap_name"`
		Cluster       string `json:"cluster"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		fmt.Printf("绑定参数失败, %v\n", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Bind绑定参数失败" + err.Error(),
			"data": nil,
		})
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
	err = service.ConfigMap.DeleteConfigMap(client, params.Namespace, params.ConfigmapName)
	if err != nil {
		fmt.Printf("删除ConfigMap失败%v\n", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success, 删除ConfigMap成功",
		"data": nil,
	})
}