package app

import (
	"github.com/gin-gonic/gin"
)

// BaseApp 基础应用实现，提供默认实现
type BaseApp struct {
	name        string
	initialized bool
	booted      bool
}

// NewBaseApp 创建基础应用
func NewBaseApp(name string) *BaseApp {
	return &BaseApp{
		name:        name,
		initialized: false,
		booted:      false,
	}
}

// Name 应用名称
func (a *BaseApp) Name() string {
	return a.name
}

// Initialize 初始化应用
func (a *BaseApp) Initialize() error {
	// 默认实现，子类可重写
	a.initialized = true
	return nil
}

// RegisterRoutes 注册路由
func (a *BaseApp) RegisterRoutes(group *gin.RouterGroup) {
	// 默认实现，子类必须重写
}

// Middlewares 获取应用中间件
func (a *BaseApp) Middlewares() []gin.HandlerFunc {
	// 默认实现，子类可重写
	return []gin.HandlerFunc{}
}

// Boot 启动应用
func (a *BaseApp) Boot() error {
	// 默认实现，子类可重写
	a.booted = true
	return nil
}

// Shutdown 关闭应用
func (a *BaseApp) Shutdown() error {
	// 默认实现，子类可重写
	return nil
}

// IsInitialized 检查应用是否已初始化
func (a *BaseApp) IsInitialized() bool {
	return a.initialized
}

// IsBooted 检查应用是否已启动
func (a *BaseApp) IsBooted() bool {
	return a.booted
}
