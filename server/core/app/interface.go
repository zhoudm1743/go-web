// Package app provides application management functionality
package app

import (
	"github.com/gin-gonic/gin"
)

// AppInterface 应用接口定义
type AppInterface interface {
	// Name 应用名称
	Name() string

	// Initialize 初始化应用
	Initialize() error

	// RegisterRoutes 注册路由
	RegisterRoutes(group *gin.RouterGroup)

	// Middlewares 获取应用中间件
	Middlewares() []gin.HandlerFunc

	// Boot 启动应用
	Boot() error

	// Shutdown 关闭应用
	Shutdown() error
}

// AppManager 应用管理器接口
type AppManager interface {
	// RegisterApp 注册应用
	RegisterApp(app AppInterface) error

	// GetApp 获取应用
	GetApp(name string) (AppInterface, bool)

	// GetApps 获取所有应用
	GetApps() []AppInterface

	// InitializeApps 初始化所有应用
	InitializeApps() error

	// BootApps 启动所有应用
	BootApps() error

	// ShutdownApps 关闭所有应用
	ShutdownApps() error
}
