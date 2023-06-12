package middle

import (
	"github.com/gin-gonic/gin"
	"k8sManagerApi/utils"
	"net/http"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 对登录接口放行
		if len(c.Request.URL.String()) >= 10 && c.Request.URL.String()[0:10] == "/api/login" {
			c.Next()
		} else if c.Request.URL.String()[0:10] == "/download/" { // 对下载接口放行
			c.Next()
		} else {
			// 获取Header中的Authorization
			token := c.Request.Header.Get("Authorization")
			if token == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"code": http.StatusBadRequest,
					"msg":  "未登录，无权限访问",
					"data": nil,
				})
				c.Abort()
				return
			}
			// 解析Token
			claims, err := utils.JWTToken.ParseToken(token)
			if err != nil {
				// token到期
				if err.Error() == "TokenExpired" {
					c.JSON(http.StatusUnauthorized, gin.H{
						"code": http.StatusUnauthorized,
						"msg":  "授权已过期，请重新登录",
						"data": nil,
					})
					c.Abort()
					return
				}
				// 其他解析错误
				c.JSON(http.StatusBadRequest, gin.H{
					"code": http.StatusBadRequest,
					"msg":  err.Error(),
					"data": nil,
				})
				c.Abort()
				return
			}

			// 继续交由下一个路由处理,并将解析出的信息传递下去
			c.Set("claims", claims)

			// 登录成功
			c.Next()
		}
	}
}