package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-web/apps/admin/controllers"
	"github.com/zhoudm1743/go-web/apps/admin/middlewares"
)

// InitRoutes 初始化路由
func InitRoutes(r *gin.Engine) {
	// 初始化控制器

	// 公开路由
	publicRoutes := r.Group("/admin")
	{
		// 路由示例
		// publicRoutes.GET("/example", exampleController.Example)
	}

	// 私有路由
	privateRoutes := r.Group("/admin")
	// privateRoutes.Use(middlewares.AuthMiddleware())
	{
		// 路由示例
		// privateRoutes.GET("/example", exampleController.Example)
	}
}
