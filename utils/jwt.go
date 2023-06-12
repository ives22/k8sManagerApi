package utils

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var JWTToken jwtToken

type jwtToken struct{}

// CustomClaims 定义token中包含的自定义信息以及jwt签名信息
type CustomClaims struct {
	Uid      uint   `json:"uid"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// 加解密因子
const (
	SECRET = "golang"
	//TokenExpireDuration = time.Second * 10
	TokenExpireDuration = time.Hour * 1 // 过期时间
)

// GenerateToken 生成token
func (*jwtToken) GenerateToken(uid uint, username string) (token string, err error) {
	cla := CustomClaims{
		uid,
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			IssuedAt:  time.Now().Unix(),                          // 签发时间
			Subject:   "Token",                                    // 主题
			Issuer:    "ives",                                     // 签发人
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, cla)
	token, err = tokenClaims.SignedString([]byte(SECRET))
	return token, err
}

// ParseToken 解析token
func (*jwtToken) ParseToken(tokenString string) (claims *CustomClaims, err error) {
	//使用jwt.ParseWithClaims方法解析token，这个token是前端传给我们的,获得一个*Token类型的对象
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET), nil
	})

	if err != nil {
		fmt.Printf("parse token failed, %v\n", err)
		//处理token解析后的各种错误
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("TokenMalformed")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("TokenExpired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("TokenNotValidYet")
			} else {
				return nil, errors.New("TokenInvalid")
			}
		}
	}

	//转换成*CustomClaims类型并返回
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("解析Token失败")
}