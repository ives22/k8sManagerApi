package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8sManagerApi/service"
	"net/http"
)

var StatefulSet statefulSet

type statefulSet struct{}

// GetStatefulSetsHandler 获取StatefulSet列表，支持分页、过滤、排序
func (s *statefulSet) GetStatefulSetsHandler(ctx *gin.Context) {
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
	data, err := service.StatefulSet.GetStatefulSets(client, params.FilterName, params.Namespace, params.Limit, params.Page)
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
		"msg":  "success, 获取StatefulSet列表成功",
		"data": data,
	})
}

// GetStatefulSetDetailHandler 获取StatefulSet详情
func (s *statefulSet) GetStatefulSetDetailHandler(ctx *gin.Context) {
	//    处理请求参数
	// 匿名结构体，用于定义传入参数，get请求为form格式，其它请求为json格式
	params := new(struct {
		StatefulSetName string `form:"statefulset_name"`
		Namespace       string `form:"namespace"`
		Cluster         string `form:"cluster"`
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
	data, err := service.StatefulSet.GetStatefulSetDetail(client, params.StatefulSetName, params.Namespace)
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
		"msg":  "success, 获取StatefulSet详情成功",
		"data": data,
	})
}

// DeleteStatefulSetHandler 删除StatefulSet
func (s *statefulSet) DeleteStatefulSetHandler(ctx *gin.Context) {
	params := new(struct {
		StatefulSetName string `json:"statefulset_name"`
		Namespace       string `json:"namespace"`
		Cluster         string `json:"cluster"`
	})
	// form格式使用Bind方法，json格式使用SholdBindJson方法
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
	// 调用Service层方法进行删除
	if err := service.StatefulSet.DeleteStatefulSet(client, params.StatefulSetName, params.Namespace); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success, 删除StatefulSet成功",
		"data": nil,
	})
}

// UpdateStatefulSetHandler 更新StatefulSet
func (s *statefulSet) UpdateStatefulSetHandler(ctx *gin.Context) {
	params := new(struct {
		Namespace string `json:"namespace"`
		Content   string `json:"content"`
		Cluster   string `json:"cluster"`
	})
	// form格式使用Bind方法，json格式使用SholdBindJson方法
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
	err = service.StatefulSet.UpdateStatefulSet(client, params.Namespace, params.Content)
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
		"msg":  "success, 更新StatefulSet成功",
		"data": nil,
	})
}

// GetStatefulSetNumPerNpHandler 获取每个Namespace的StatefulSet的数量
func (s *statefulSet) GetStatefulSetNumPerNpHandler(ctx *gin.Context) {
	params := new(struct {
		Cluster string `form:"cluster"`
	})
	if err := ctx.Bind(params); err != nil {
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
	data, err := service.StatefulSet.GetStatefulSetNumPerNp(client)
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
		"msg":  "success, 获取每个namespace的StatefulSet数量成功",
		"data": data,
	})
}