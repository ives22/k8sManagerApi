package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8sManagerApi/service"
	"net/http"
)

var Auth auth

type auth struct{}

func (a *auth) LoginHandler(ctx *gin.Context) {
	params := new(struct {
		Username string `json:"username"`
		Password string `json:"password"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		fmt.Printf("绑定参数失败, %v\n", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "参数传递错误, " + err.Error(),
			"data": nil,
		})
		return
	}
	token, err := service.Auth.Login(params.Username, params.Password)
	if err != nil {
		fmt.Printf("登录失败, %v\n", err.Error())
		ctx.JSON(http.StatusForbidden, gin.H{
			"code": http.StatusForbidden,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success, certification passed",
		"data": gin.H{
			"token": token,
		},
	})
}