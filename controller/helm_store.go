package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"k8sManagerApi/model"
	"k8sManagerApi/service"
	"net/http"
)

var HelmStore helmStore

type helmStore struct{}

// ListReleasesHandler 获取已安装的Release列表
func (h *helmStore) ListReleasesHandler(ctx *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
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
	actionConfig, err := service.HelmConfig.GetAc(params.Cluster, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.HelmStore.ListRelease(actionConfig, params.FilterName)
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
		"msg":  "success, 获取Release列表成功",
		"data": data,
	})
}

// DetailReleaseHandler 获取Release详情
func (h *helmStore) DetailReleaseHandler(ctx *gin.Context) {
	params := new(struct {
		Release   string `form:"release"`
		Namespace string `form:"namespace"`
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
	actionConfig, err := service.HelmConfig.GetAc(params.Cluster, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.HelmStore.DetailRelease(actionConfig, params.Release)
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
		"msg":  "success, 获取Release详情成功",
		"data": data,
	})
}

// InstallReleaseHandler 安装Release
func (h *helmStore) InstallReleaseHandler(ctx *gin.Context) {
	params := new(struct {
		Release   string `json:"release"`
		Chart     string `json:"chart"`
		Namespace string `json:"namespace"`
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
	actionConfig, err := service.HelmConfig.GetAc(params.Cluster, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err = service.HelmStore.InstallRelease(actionConfig, params.Release, params.Chart, params.Namespace)
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
		"msg":  "success, 安装Release成功",
		"data": nil,
	})
}

// UninstallReleaseHandler 卸载Release
func (h *helmStore) UninstallReleaseHandler(ctx *gin.Context) {
	params := new(struct {
		Release   string `json:"release"`
		Namespace string `json:"namespace"`
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
	actionConfig, err := service.HelmConfig.GetAc(params.Cluster, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err = service.HelmStore.UninstallRelease(actionConfig, params.Release, params.Namespace)
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
		"msg":  "success, 卸载Release成功",
		"data": nil,
	})
}

// UploadChartFileHandler Chart文件上传
func (h *helmStore) UploadChartFileHandler(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("chart")
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取上传信息失败, %v", err.Error()))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "获取上传信息失败" + err.Error(),
			"data": nil,
		})
		return
	}
	err = service.HelmStore.UploadChartFile(file, header)
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
		"msg":  "success, 上传Chart文件成功",
		"data": nil,
	})
}

// DeleteChartFileHandler Chart文件删除
func (h *helmStore) DeleteChartFileHandler(ctx *gin.Context) {
	params := new(struct {
		Chart string `json:"chart"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		zap.L().Error(fmt.Sprintf("Bind绑定参数失败, %v", err.Error()))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "Bind绑定参数失败" + err.Error(),
			"data": nil,
		})
		return
	}
	err := service.HelmStore.DeleteChartFile(params.Chart)
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
		"msg":  "success, 删除Chart文件成功",
		"data": nil,
	})
}

// GetChartsHandler 获取chart列表
func (h *helmStore) GetChartsHandler(ctx *gin.Context) {
	params := new(struct {
		Name  string `form:"name"`
		Page  int    `form:"page"`
		Limit int    `form:"limit"`
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
	data, err := service.HelmStore.GetCharts(params.Name, params.Page, params.Limit)
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
		"msg":  "success, 获取Chart列表成功",
		"data": data,
	})
}

// AddChartHandler 新增Chart
func (h *helmStore) AddChartHandler(ctx *gin.Context) {
	params := new(model.Chart)
	if err := ctx.ShouldBindJSON(params); err != nil {
		zap.L().Error(fmt.Sprintf("Bind绑定参数失败, %v", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Bind绑定参数失败" + err.Error(),
			"data": nil,
		})
		return
	}
	err := service.HelmStore.AddChart(params)
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
		"msg":  "success, 新增Chart成功",
		"data": nil,
	})
}

// UpdateChartHandler 更新Chart
func (h *helmStore) UpdateChartHandler(ctx *gin.Context) {
	params := new(model.Chart)
	if err := ctx.ShouldBindJSON(params); err != nil {
		zap.L().Error(fmt.Sprintf("Bind绑定参数失败, %v", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Bind绑定参数失败" + err.Error(),
			"data": nil,
		})
		return
	}
	err := service.HelmStore.UpdateChart(params)
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
		"msg":  "success, 更新Chart成功",
		"data": nil,
	})
}

// DeleteChartHandler 删除Chart
func (h *helmStore) DeleteChartHandler(ctx *gin.Context) {
	params := new(model.Chart)
	if err := ctx.ShouldBindJSON(params); err != nil {
		zap.L().Error(fmt.Sprintf("Bind绑定参数失败, %v", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Bind绑定参数失败" + err.Error(),
			"data": nil,
		})
		return
	}
	err := service.HelmStore.DeleteChart(params)
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
		"msg":  "success, 删除Chart成功",
		"data": nil,
	})
}