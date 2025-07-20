package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-web/core/response"
	"github.com/zhoudm1743/go-web/core/utils"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取JWT声明
		claims, err := utils.GetClaims(c)
		if err != nil {
			response.Fail(c, response.TokenInvalid)
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中，方便后续使用
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roleID", claims.RoleID)

		c.Next()
	}
}
