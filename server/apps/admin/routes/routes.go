package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-web/apps/admin/controllers"
	"github.com/zhoudm1743/go-web/apps/admin/middlewares"
)

// InitRoutes 初始化路由
func InitRoutes(r *gin.RouterGroup) {
	// 初始化控制器
	authController := controllers.NewAuthController()
	adminController := controllers.NewAdminController()
	menuController := controllers.NewMenuController()
	roleController := controllers.NewRoleController()
	codeGenController := controllers.NewCodeGenController()

	// 公开路由
	publicRoutes := r.Group("/admin")
	{
		// 认证相关路由
		publicRoutes.POST("/login", authController.Login)
		publicRoutes.POST("/refresh", authController.RefreshToken)
	}

	// 私有路由
	privateRoutes := r.Group("/admin")
	privateRoutes.Use(middlewares.AuthMiddleware())
	{
		// 认证相关路由
		privateRoutes.GET("/me", authController.GetUserInfo)
		privateRoutes.POST("/logout", authController.Logout)

		// 管理员路由
		privateRoutes.GET("/admins", adminController.GetAdmins)
		privateRoutes.POST("/admin", adminController.CreateAdmin)
		privateRoutes.PUT("/admin", adminController.UpdateAdmin)
		privateRoutes.DELETE("/admin/:id", adminController.DeleteAdmin)

		// 菜单路由
		privateRoutes.GET("/menus", menuController.GetMenus)
		privateRoutes.GET("/menus/tree", menuController.GetMenuTree)
		privateRoutes.POST("/menu", menuController.CreateMenu)
		privateRoutes.PUT("/menu", menuController.UpdateMenu)
		privateRoutes.DELETE("/menu/:id", menuController.DeleteMenu)

		// 角色路由
		privateRoutes.GET("/roles", roleController.GetRoles)
		privateRoutes.POST("/role", roleController.CreateRole)
		privateRoutes.PUT("/role", roleController.UpdateRole)
		privateRoutes.DELETE("/role/:id", roleController.DeleteRole)
		privateRoutes.PUT("/role/menu", roleController.AssignMenu)

		// 代码生成器路由
		codeGenController.RegisterRoutes(privateRoutes)
	}
}
