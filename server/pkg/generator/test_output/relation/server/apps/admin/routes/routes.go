package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-web/apps/admin/controllers"
)

// InitRoutes 初始化路由
func InitRoutes(r *gin.Engine) {
	// 初始化控制器
	articleController := controllers.NewArticleController()
	articleController := controllers.NewArticleController()
	// 私有路由
	privateRoutes := r.Group("/admin")
	{
		// 路由示例
		// privateRoutes.GET("/example", exampleController.Example)
	}
}