package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterGlobalMiddlewares 注册全局中间件
func RegisterGlobalMiddlewares(router *gin.Engine) {
	// 跨域中间件
	router.Use(corsMiddleware())

	// 请求日志中间件在app.go中已经添加
}

// RegisterGlobalRoutes 注册全局路由
func RegisterGlobalRoutes(router *gin.Engine) {
	// 健康检查路由
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "服务运行正常",
		})
	})

	// 版本信息路由
	router.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version": "1.0.0",
			"name":    "Go-Web Framework",
		})
	})

	// 添加其他全局路由...
}

// corsMiddleware 跨域中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// 处理OPTIONS请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
