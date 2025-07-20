package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-web/core/app"
)

// APIApp API应用
type APIApp struct {
	*app.BaseApp
}

// NewAPIApp 创建API应用
func NewAPIApp() *APIApp {
	return &APIApp{
		BaseApp: app.NewBaseApp("api"),
	}
}

// Initialize 初始化应用
func (a *APIApp) Initialize() error {
	// 调用父类初始化
	if err := a.BaseApp.Initialize(); err != nil {
		return err
	}

	// 避免使用可能未初始化的facades
	// 使用fmt.Println代替facades.Log
	fmt.Println("初始化API应用:", a.Name())
	return nil
}

// RegisterRoutes 注册路由
func (a *APIApp) RegisterRoutes(group *gin.RouterGroup) {
	// 注册API路由
	group.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"app":     "api",
			"message": "API服务",
			"version": "v1",
		})
	})

	// 用户API
	userGroup := group.Group("/users")
	userGroup.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"users": []gin.H{
				{"id": 1, "name": "用户1"},
				{"id": 2, "name": "用户2"},
			},
		})
	})

	// 数据API
	dataGroup := group.Group("/data")
	dataGroup.GET("/stats", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"stats": gin.H{
				"userCount":    100,
				"articleCount": 58,
			},
		})
	})
}

// Boot 启动应用
func (a *APIApp) Boot() error {
	if err := a.BaseApp.Boot(); err != nil {
		return err
	}

	// 避免使用可能未初始化的facades
	fmt.Println("启动API应用:", a.Name())
	return nil
}

// Middlewares 获取应用中间件
func (a *APIApp) Middlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		func(c *gin.Context) {
			// 避免使用可能未初始化的facades
			fmt.Println("API中间件被调用")
			c.Next()
		},
	}
}
