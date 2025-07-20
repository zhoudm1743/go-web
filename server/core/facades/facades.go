package facades

import (
	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-web/core/app"
	"github.com/zhoudm1743/go-web/core/cache"
	"github.com/zhoudm1743/go-web/core/log"
	"gorm.io/gorm"
)

// 全局变量
var (
	appManagerInstance AppManagerInterface
	loggerInstance     log.Logger
	dbInstance         *gorm.DB
	routeInstance      RouterInterface
	cacheInstance      cache.Cache
)

// AppManagerInterface 应用管理器接口
type AppManagerInterface interface {
	// RegisterApp 注册应用
	RegisterApp(app app.AppInterface) error
	// GetApp 获取应用
	GetApp(name string) (app.AppInterface, bool)
	// GetApps 获取所有应用
	GetApps() []app.AppInterface
	// InitializeApps 初始化所有应用
	InitializeApps() error
	// BootApps 启动所有应用
	BootApps() error
	// ShutdownApps 关闭所有应用
	ShutdownApps() error
}

// RouterGroup 路由组接口 - 简化，直接使用Gin的RouterGroup
type RouterGroup = *gin.RouterGroup

// RouterInterface 路由器接口 - 简化
type RouterInterface interface {
	// Group 创建路由组
	Group(path string, handlers ...interface{}) RouterGroup
	// Run 运行HTTP服务
	Run() error
	// Shutdown 关闭HTTP服务
	Shutdown() error
}

// SetAppManager 设置全局AppManager
func SetAppManager(manager AppManagerInterface) {
	appManagerInstance = manager
}

// SetLog 设置全局Logger
func SetLog(l log.Logger) {
	loggerInstance = l
}

// SetDB 设置全局DB
func SetDB(database *gorm.DB) {
	dbInstance = database
}

// SetRoute 设置全局Router
func SetRoute(r RouterInterface) {
	routeInstance = r
}

// SetCache 设置全局Cache
func SetCache(c cache.Cache) {
	cacheInstance = c
}

// AppManager 获取全局AppManager
func AppManager() AppManagerInterface {
	return appManagerInstance
}

// Log 获取全局Logger
func Log() log.Logger {
	return loggerInstance
}

// DB 获取全局DB
func DB() *gorm.DB {
	return dbInstance
}

// Route 获取全局Router
func Route() RouterInterface {
	return routeInstance
}

// Cache 获取全局Cache
func Cache() cache.Cache {
	return cacheInstance
}
