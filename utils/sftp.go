package utils

import (
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"path/filepath"
	"time"
)

//var Cli cli

type Cli struct {
	Username string
	Password string
	Host     string
	Port     int
	Client   *ssh.Client
}

// Connection 建立SSH连接，初始化 client
func (c *Cli) Connection() (err error) {
	// 远程服务器的SSH连接配置
	config := &ssh.ClientConfig{
		User: c.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 10,
	}
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	// 连接远程服务器
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to connect to the server, %v", err))
		return err
	}
	c.Client = sshClient
	return nil
}

// ScpToServer 本地拷贝文件到服务服务器上
func (c *Cli) ScpToServer(srcFile, destFile string) (err error) {
	if c.Client == nil {
		if err := c.Connection(); err != nil {
			return err
		}
	}

	// 1 基于ssh client，创建 sftp 客户端
	sftpClient, err := sftp.NewClient(c.Client)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to init sftp client, %v", err))
		return err
	}
	defer sftpClient.Close()

	// 2 检查远程文件是否存在，如果存在则跳过
	_, err = sftpClient.Stat(destFile)
	if err == nil {
		zap.L().Warn(fmt.Sprintf("File %s already exists on the server. Skipping...", destFile))
		return errors.New(fmt.Sprintf("File %s already exists on the server. Skipping", destFile))
	}

	// 3 判断远程上传目录是否存在，如果不存在则创建
	dir := filepath.Dir(destFile)
	file := filepath.Base(destFile)
	_, err = sftpClient.Stat(dir)
	if os.IsNotExist(err) {
		sftpClient.MkdirAll(dir)
	}

	// 4 创建远程服务器文件
	remoteFile, err := sftpClient.Create(filepath.Join(dir, file))
	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to create remote file, %v", err))
		return err
	}
	defer remoteFile.Close()

	// 5 打开本地文件
	localFile, err := os.Open(srcFile)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to open local file, %v", err))
		return err
	}
	defer localFile.Close()

	// 6 将本地文件内容复制到远程文件
	_, err = io.Copy(remoteFile, localFile)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to copy local file to remote file, %v", err))
		return err
	}

	return nil
}