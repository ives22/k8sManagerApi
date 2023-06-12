package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8sManagerApi/service"
	"net/http"
)

var AllRes allRes

type allRes struct{}

// GetAllNumHandler 获取集群的所有资源
func (a *allRes) GetAllNumHandler(ctx *gin.Context) {
	params := new(struct {
		Cluster string `form:"cluster"`
	})
	if err := ctx.Bind(params); err != nil {
		fmt.Printf("绑定参数失败, %v\n", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  errors.New("绑定参数失败," + err.Error()),
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
	data, errs := service.AllRes.GetAllNum(client)
	if len(errs) > 0 {
		fmt.Printf("绑定参数失败, %v\n", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  errs,
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success, 获取集群资源数量成功",
		"data": data,
	})
}