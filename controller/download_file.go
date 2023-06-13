package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8sManagerApi/config"
	"net/http"
	"os"
	"path/filepath"
)

var Download download

type download struct{}

// DownloadHandler 用于服务器下载chart文件
func (d *download) DownloadHandler(ctx *gin.Context) {
	fileName := ctx.Param("filename") // 从URL中获取文件名
	fmt.Printf("文件名：%v\n", fileName)
	// 拼接文件路径
	filePath := filepath.Join(config.Conf.UploadPath, fileName)
	// 判断文件是否存在
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
			"msg":  "no such file",
			"data": nil,
		})
		return
	}

	// 设置响应头，指定下载文件的名称
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))

	// 设置响应内容类型为二进制流
	ctx.Header("Content-Type", "application/octet-stream")

	// 发送文件内容给客户端
	ctx.File(filePath)
}