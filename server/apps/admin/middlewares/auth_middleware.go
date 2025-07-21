package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-web/apps/admin/models"
	"github.com/zhoudm1743/go-web/core/facades"
	"github.com/zhoudm1743/go-web/core/middleware"
	"github.com/zhoudm1743/go-web/core/response"
	"github.com/zhoudm1743/go-web/core/utils"
)

// PermissionAuth 权限认证中间件
func PermissionAuth() gin.HandlerFunc {
	return middleware.CasbinHandler()
}

// AdminAuth 管理员认证中间件
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户信息
		claims, err := utils.GetClaims(c)
		if err != nil {
			response.Fail(c, response.TokenInvalid)
			c.Abort()
			return
		}

		// 查询用户角色
		var admin models.Admin
		if err := facades.DB().First(&admin, claims.UserID).Error; err != nil {
			response.Fail(c, response.NoPermission)
			c.Abort()
			return
		}

		// 检查角色
		roles := admin.GetRoles()
		isAdmin := false
		for _, role := range roles {
			if role == "super" { // 修改为检查 "super" 角色
				isAdmin = true
				break
			}
		}

		if !isAdmin {
			response.Fail(c, response.NoPermission)
			c.Abort()
			return
		}

		c.Next()
	}
}
