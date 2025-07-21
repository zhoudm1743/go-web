package admin

import (
	"github.com/gin-gonic/gin"
)

// Register 注册应用
func Register(r *gin.Engine) {
	// 初始化路由
	InitRoutes(r)
}
