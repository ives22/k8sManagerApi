package service

import (
	"k8sManagerApi/dao"
	"k8sManagerApi/model"
	"k8sManagerApi/utils"
)

var User user

type user struct{}

type UserCreate struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// CreateUser 创建用户
func (u *user) CreateUser(user *UserCreate) (err error) {
	// 对密码进行加密处理
	md5Password := utils.ToMd5(user.Password)
	newUser := &model.User{
		Username: user.Username,
		Password: md5Password,
	}
	dao.User.AddUser(newUser)
	if err != nil {
		return err
	}
	return nil
}

// GetUser 查询用户
func (u *user) GetUser(username string) (data *model.User, err error) {
	data, err = dao.User.GetUserByName(username)
	if err != nil {
		return nil, err
	}
	return data, err
}