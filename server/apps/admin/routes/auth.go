package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-web/apps/admin/controllers"
	"github.com/zhoudm1743/go-web/core/middleware"
)

// RegisterAuthRoutes 注册认证相关路由
func RegisterAuthRoutes(r *gin.RouterGroup) {
	// 创建控制器
	authController := controllers.NewAuthController()

	// 公开接口 - 无需认证
	publicRouter := r.Group("/auth")
	{
		publicRouter.POST("/login", authController.Login)   // 登录
		publicRouter.POST("/logout", authController.Logout) // 登出
	}

	// 受保护接口 - 需要认证
	authRouter := r.Group("")
	authRouter.Use(middleware.JWTAuth())
	{
		authRouter.GET("/user/info", authController.GetUserInfo)     // 获取用户信息
		authRouter.GET("/auth/codes", authController.GetAccessCodes) // 获取权限码
		authRouter.GET("/menu/all", authController.GetUserRoutes)    // 获取用户菜单
	}
}
