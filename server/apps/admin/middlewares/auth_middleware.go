package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-web/core/response"
	"github.com/zhoudm1743/go-web/core/utils"
)

// JWTAuthMiddleware JWT认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从JWT中获取声明
		claims, err := utils.GetClaims(c)
		if err != nil {
			response.Fail(c, response.TokenInvalid)
			c.Abort()
			return
		}

		// 将用户信息存储在上下文中
		c.Set("claims", claims)
		c.Set("userID", claims.UserID)
		c.Set("uuid", claims.UUID)
		c.Set("roles", claims.Roles)
		c.Set("roleID", claims.RoleID)

		c.Next()
	}
}

// RoleAuthMiddleware 角色认证中间件
func RoleAuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("roles")
		if !exists {
			response.Fail(c, response.NoPermission)
			c.Abort()
			return
		}

		// 检查用户角色是否在允许的角色列表中
		userRoles := roles.([]string)
		allowed := false
		for _, userRole := range userRoles {
			for _, allowedRole := range allowedRoles {
				if allowedRole == userRole {
					allowed = true
					break
				}
			}
			if allowed {
				break
			}
		}

		if !allowed {
			response.Fail(c, response.NoPermission)
			c.Abort()
			return
		}

		c.Next()
	}
}
