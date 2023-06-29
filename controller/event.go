package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"k8sManagerApi/service"
	"net/http"
)

var Event event

type event struct{}

// GetEventsHandler 获取事件列表
func (e *event) GetEventsHandler(ctx *gin.Context) {
	params := new(struct {
		Name    string `form:"name"`
		Cluster string `form:"cluster"`
		Page    int    `form:"page"`
		Limit   int    `form:"limit"`
	})
	if err := ctx.Bind(params); err != nil {
		zap.L().Error(fmt.Sprintf("Bind绑定参数失败, %v", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  errors.New("绑定参数失败, " + err.Error()),
			"data": nil,
		})
		return
	}
	data, err := service.Event.GetEvents(params.Name, params.Cluster, params.Page, params.Limit)
	if err != nil {
		zap.L().Error(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "获取Evens列表成功",
		"data": data,
	})
}