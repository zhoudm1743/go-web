package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-web/apps/admin/controllers"
	"github.com/zhoudm1743/go-web/apps/admin/middlewares"
	"github.com/zhoudm1743/go-web/core/middleware"
)

// RegisterRoutes 注册路由
func RegisterRoutes(r *gin.RouterGroup) {
	// 初始化控制器
	authController := controllers.NewAuthController()
	userController := controllers.NewAdminController()
	menuController := controllers.NewMenuController()
	roleController := controllers.NewRoleController()

	// 公共路由（无需认证）
	{
		// 登录认证
		r.POST("/login", authController.Login)
	}

	// 需要认证的路由
	privateRoutes := r.Group("")
	privateRoutes.Use(middleware.JWTAuth())
	{
		// 认证相关
		privateRoutes.POST("/logout", authController.Logout)
		privateRoutes.GET("/userInfo", authController.GetUserInfo)
		privateRoutes.GET("/accessCodes", authController.GetAccessCodes)

		// 用户路由菜单
		privateRoutes.GET("/getUserRoutes", authController.GetUserRoutes)

		// 所有菜单（用于菜单管理页面）
		privateRoutes.GET("/getAllRoutes", authController.GetAllRoutes)

		// 角色相关
		roleGroup := privateRoutes.Group("/role")
		{
			// 角色列表 - 简单列表，用于下拉选择
			roleGroup.GET("/list", roleController.GetRoleList)
			roleGroup.POST("/create", roleController.CreateRole)
			roleGroup.PUT("/update", roleController.UpdateRole)
			roleGroup.DELETE("/delete/:id", roleController.DeleteRole)
			roleGroup.GET("/menus", roleController.GetRoleMenus)
			roleGroup.POST("/menus", roleController.UpdateRoleMenus)
		}

		// 用户管理
		userGroup := privateRoutes.Group("/user")
		userGroup.Use(middlewares.PermissionAuth())
		{
			userGroup.GET("/list", userController.GetAdmins)
			userGroup.POST("/create", userController.CreateAdmin)
			userGroup.PUT("/update", userController.UpdateAdmin)
			userGroup.DELETE("/delete/:id", userController.DeleteAdmin)
		}
	}

	// 需要认证和管理员权限的路由
	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middleware.JWTAuth())
	adminRoutes.Use(middlewares.AdminAuth())
	{
		// 菜单管理
		menuGroup := adminRoutes.Group("/menu")
		{
			menuGroup.GET("/list", menuController.GetMenus)
			menuGroup.POST("/create", menuController.CreateMenu)
			menuGroup.PUT("/update", menuController.UpdateMenu)
			menuGroup.DELETE("/delete/:id", menuController.DeleteMenu)
		}
	}
}
