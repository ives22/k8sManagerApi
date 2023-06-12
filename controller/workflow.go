package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8sManagerApi/service"
	"net/http"
)

var Workflow workflow

type workflow struct{}

// GetWorkflowsHandler 获取列表分页查询
func (w *workflow) GetWorkflowsHandler(ctx *gin.Context) {
	params := new(struct {
		Name      string `form:"name"`
		Namespace string `form:"namespace"`
		Page      int    `form:"page"`
		Limit     int    `form:"limit"`
		Cluster   string `form:"cluster"`
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

	data, err := service.Workflow.GetWorkflows(params.Name, params.Namespace, params.Cluster, params.Limit, params.Page)
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
		"msg":  "success, 获取workflow列表成功",
		"data": data,
	})
}

// GetWorkflowDetailHandler 查询workflow单条数据
func (w *workflow) GetWorkflowDetailHandler(ctx *gin.Context) {
	params := new(struct {
		ID int `form:"id"`
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

	data, err := service.Workflow.GetWorkflowDetail(params.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "查询Workflow单条数据失败" + err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success, 查询Workflow单条数据成功",
		"data": data,
	})
}

// CreateWorkflowHandler 创建workflow
func (w *workflow) CreateWorkflowHandler(ctx *gin.Context) {
	var (
		wc  = &service.WorkflowCreate{}
		err error
	)
	if err := ctx.ShouldBindJSON(wc); err != nil {
		fmt.Printf("绑定参数失败, %v\n", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Bind绑定参数失败" + err.Error(),
			"data": nil,
		})
		return
	}

	fmt.Println("controller", 999)
	client, err := service.K8s.GetClient(wc.Cluster)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	if err = service.Workflow.CreateWorkflow(client, wc); err != nil {
		fmt.Printf("创建Workflow失败, %v\n", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success, 创建Workflow成功",
		"data": nil,
	})
}

// DeleteWorkflowHandler 删除workflow
func (w *workflow) DeleteWorkflowHandler(ctx *gin.Context) {
	params := new(struct {
		ID      int    `json:"id"`
		Cluster string `json:"cluster"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		fmt.Printf("Bind绑定参数失败, %v\n", err.Error())
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
	if err := service.Workflow.DeleteWorkflow(client, params.ID); err != nil {
		fmt.Printf("删除Workflow失败, %v\n", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success, 删除Workflow成功",
		"data": nil,
	})
}