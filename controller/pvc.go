package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"k8sManagerApi/service"
	"net/http"
)

var Pvc pvc

type pvc struct{}

// GetPvcsHandler 获取PVC列表
func (p *pvc) GetPvcsHandler(ctx *gin.Context) {
	// 处理请求参数
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
		Cluster    string `form:"cluster"`
	})
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
	data, err := service.Pvc.GetPvcs(client, params.Namespace, params.FilterName, params.Limit, params.Page)
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
		"msg":  "success, 获取PVC列表成功",
		"data": data,
	})
}

// GetPvcDetailHandler 获取PVC详情
func (p *pvc) GetPvcDetailHandler(ctx *gin.Context) {
	params := new(struct {
		Namespace string `form:"namespace"`
		PvcName   string `form:"pvc_name"`
		Cluster   string `form:"cluster"`
	})
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
	data, err := service.Pvc.GetPvcDetail(client, params.Namespace, params.PvcName)
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
		"msg":  "success, 获取PVC详情成功",
		"data": data,
	})
}

// UpdatePvcHandler 更新PVC
func (p *pvc) UpdatePvcHandler(ctx *gin.Context) {
	params := new(struct {
		Namespace string `json:"namespace"`
		Content   string `json:"Content"`
		Cluster   string `json:"cluster"`
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
	err = service.Pvc.UpdatePvc(client, params.Namespace, params.Content)
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
		"msg":  "success, 更新PVC成功",
		"data": nil,
	})
}

// DeletePvcHandler 删除PVC
func (p *pvc) DeletePvcHandler(ctx *gin.Context) {
	params := new(struct {
		Namespace string `json:"namespace"`
		PvcName   string `json:"pvc_name"`
		Cluster   string `json:"cluster"`
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
	err = service.Pvc.DeletePvc(client, params.Namespace, params.PvcName)
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
		"msg":  "success, 删除PVC成功",
		"data": nil,
	})
}