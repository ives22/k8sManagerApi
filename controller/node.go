package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8sManagerApi/service"
	"net/http"
)

var Node node

type node struct{}

// GetNodesHandler GetNodes 获取node列表
func (n *node) GetNodesHandler(ctx *gin.Context) {
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
	data, err := service.Node.GetNodes(client, params.FilterName, params.Limit, params.Page)
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
		"msg":  "success, 获取Node列表成功",
		"data": data,
	})
}

// GetNodeDetailHandler 获取node详情
func (n *node) GetNodeDetailHandler(ctx *gin.Context) {
	// 匿名结构体，用于定义传入参数，get请求为form格式，其它请求为json格式
	params := new(struct {
		NodeName string `form:"node_name"`
		Cluster  string `form:"cluster"`
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
	data, err := service.Node.GetNodeDetail(client, params.NodeName)
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
		"msg":  "success, 获取Node详情成功",
		"data": data,
	})
}