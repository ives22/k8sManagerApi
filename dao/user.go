package dao

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"k8sManagerApi/db/mysql"
	"k8sManagerApi/model"
)

var User user

type user struct{}

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// GetUserByName 根据用户名查找用户信息
func (u *user) GetUserByName(name string) (user *model.User, err error) {
	user = &model.User{}
	tx := mysql.DB.Where("username = ?", name).First(&user)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		zap.L().Error(fmt.Sprintf("查询用户失败, %v", tx.Error.Error()))
		return nil, errors.New("查询用户失败," + tx.Error.Error())
	}
	return user, nil
}

// AddUser 创建用户
func (u *user) AddUser(user *model.User) (err error) {
	tx := mysql.DB.Create(&user)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		zap.L().Error(fmt.Sprintf("创建用户失败, %v", tx.Error.Error()))
		return errors.New("创建用户失败," + tx.Error.Error())
	}
	return nil
}

// ChangePassword 修改密码
func (u *user) ChangePassword(name, password string) (err error) {
	user := &model.User{}
	_ = mysql.DB.Where("username = ?", name).First(&user)
	mysql.DB.Model(user).Select("password").Update("password", password)
	return nil
}