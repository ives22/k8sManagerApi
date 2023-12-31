package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"k8sManagerApi/service"
	"net/http"
)

var Auth auth

type auth struct{}

// LoginHandler 用户登录
func (a *auth) LoginHandler(ctx *gin.Context) {
	params := new(struct {
		Username string `json:"username"`
		Password string `json:"password"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		zap.L().Error(fmt.Sprintf("绑定参数失败, %v", err.Error()))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "参数传递错误, " + err.Error(),
			"data": nil,
		})
		return
	}
	token, err := service.Auth.Login(params.Username, params.Password)
	if err != nil {
		zap.L().Error(fmt.Sprintf("登录失败, %v", err.Error()))
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

// RegisterHandler 用户注册
func (a *auth) RegisterHandler(ctx *gin.Context) {
	//params := new(struct {
	//	Username string `json:"username"`
	//	Password string `json:"password"`
	//})
	userinfo := &service.UserCreate{}
	if err := ctx.ShouldBindJSON(&userinfo); err != nil {
		zap.L().Error(fmt.Sprintf("绑定参数失败, %v", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "绑定参数失败," + err.Error(),
			"data": nil,
		})
		return
	}
	if err := service.Auth.Register(userinfo); err != nil {
		zap.L().Error(fmt.Sprintf("create user failed, %v", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "create user failed," + err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "create user success",
		"data": nil,
	})
}

// ChangePwdHandler 修改密码
func (a *auth) ChangePwdHandler(ctx *gin.Context) {
	params := new(struct {
		Username string `json:"username"`
		OldPwd   string `json:"old_pwd"`
		NewPwd   string `json:"new_pwd"`
	})
	if err := ctx.ShouldBindJSON(&params); err != nil {
		zap.L().Error(fmt.Sprintf("绑定参数失败, %v", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "绑定参数失败," + err.Error(),
			"data": nil,
		})
		return
	}
	// 调用service层进行密码修改
	if err := service.Auth.ChangePwd(params.Username, params.OldPwd, params.NewPwd); err != nil {
		zap.L().Error(fmt.Sprintf("change user password failed, %v", err.Error()))

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "change user password success",
		"data": nil,
	})
}