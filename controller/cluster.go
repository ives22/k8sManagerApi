package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8sManagerApi/service"
	"net/http"
	"sort"
)

var Cluster cluster

type cluster struct{}

// GetClustersHandler 获取集群列表
func (c *cluster) GetClustersHandler(ctx *gin.Context) {
	list := make([]string, 0)
	for key := range service.K8s.ClientMap {
		list = append(list, key)
	}
	fmt.Printf("list:====> %T\n", list)
	sort.Strings(list)
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success, 获取集群列表成功",
		"data": list,
	})
}