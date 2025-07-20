package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-web/core/response"
	"github.com/zhoudm1743/go-web/core/utils"
)

// CasbinHandler 权限拦截中间件
func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前用户JWT中的信息
		claims, err := utils.GetClaims(c)
		if err != nil {
			response.Fail(c, response.TokenInvalid)
			c.Abort()
			return
		}

		// 获取请求的路径
		obj := c.Request.URL.Path
		// 获取请求方法
		act := c.Request.Method
		// 获取用户的角色
		sub := strconv.Itoa(int(claims.RoleID))

		// 获取casbin实例并检查权限
		e := utils.Casbin()
		success, err := e.Enforce(sub, obj, act)

		if err != nil {
			response.Fail(c, response.SystemError)
			c.Abort()
			return
		}

		// 判断策略中是否存在
		if !success {
			response.Fail(c, response.NoPermission)
			c.Abort()
			return
		}

		c.Next()
	}
}
