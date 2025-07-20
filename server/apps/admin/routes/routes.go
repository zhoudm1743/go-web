package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-web/apps/admin/controllers"
	"github.com/zhoudm1743/go-web/apps/admin/middlewares"
	"github.com/zhoudm1743/go-web/core/middleware"
)

// RegisterRoutes 注册admin应用的所有路由
func RegisterRoutes(r *gin.RouterGroup) {
	// API前缀
	apiGroup := r.Group("/api")

	// 注册认证相关路由
	RegisterAuthRoutes(apiGroup)

	// TODO: 注册其他路由
	// RegisterUserRoutes(apiGroup)
	// RegisterRoleRoutes(apiGroup)
	// RegisterMenuRoutes(apiGroup)
}

// RegisterAdminRoutes 注册管理后台路由
func RegisterAdminRoutes(r *gin.RouterGroup) {
	// 创建控制器实例
	authController := controllers.NewAuthController()

	// 认证相关路由
	authRoutes := r.Group("/auth")
	{
		// 无需认证的路由
		authRoutes.POST("/login", authController.Login)
		// authRoutes.POST("/register", authController.Register)
		// authRoutes.POST("/refresh", authController.RefreshToken)

		// 需要认证的路由
		authProtected := authRoutes.Group("")
		authProtected.Use(middlewares.JWTAuthMiddleware())
		{
			authProtected.GET("/info", authController.GetUserInfo)
			authProtected.GET("/menus", authController.GetUserMenus)
			// authProtected.GET("/routers", authController.GetUserRouters)
			// authProtected.PUT("/password", authController.ChangePassword)
			authProtected.POST("/logout", authController.Logout)
		}
	}

	// 用户管理路由组 (需要管理员权限)
	userRoutes := r.Group("/users")
	userRoutes.Use(middlewares.JWTAuthMiddleware())
	userRoutes.Use(middleware.CasbinHandler()) // 使用Casbin权限检查
	{
		// TODO: 添加用户管理相关路由
	}

	// 角色管理路由组 (需要管理员权限)
	roleRoutes := r.Group("/roles")
	roleRoutes.Use(middlewares.JWTAuthMiddleware())
	roleRoutes.Use(middleware.CasbinHandler()) // 使用Casbin权限检查
	{
		// TODO: 添加角色管理相关路由
	}

	// 菜单管理路由组 (需要管理员权限)
	menuRoutes := r.Group("/menus")
	menuRoutes.Use(middlewares.JWTAuthMiddleware())
	menuRoutes.Use(middleware.CasbinHandler()) // 使用Casbin权限检查
	{
		// TODO: 添加菜单管理相关路由
	}
}
