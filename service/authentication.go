package service

import (
	"errors"
	"k8sManagerApi/dao"
	"k8sManagerApi/utils"
)

var Auth auth

type auth struct{}

// Login 登录 验证用户密码是否正确，并返回token
func (a *auth) Login(username, password string) (token string, err error) {
	// 首先去数据库查询，是否存在
	user, err := dao.User.GetUserByName(username)
	if err != nil {
		return "", errors.New("users do not exist")
	}
	// 如果用户存在，则进行密码校验, 将传入的密码进行md5加密，然后和数据库中的密码进行比较
	if user.Password != utils.ToMd5(password) {
		return "", errors.New("authentication failed, password error")
	}
	// 用户密码校验完成后，生成token返回给用户
	token, err = utils.JWTToken.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}