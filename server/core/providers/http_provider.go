package providers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-web/core"
	"github.com/zhoudm1743/go-web/core/conf"
	"github.com/zhoudm1743/go-web/core/facades"
	"github.com/zhoudm1743/go-web/core/http"
	"github.com/zhoudm1743/go-web/core/log"
)

// HTTPProvider HTTP服务提供者
type HTTPProvider struct {
	config *conf.Config
	logger log.Logger
	server http.HTTPServer
}

// NewHTTPProvider 创建HTTP服务提供者
func NewHTTPProvider() *HTTPProvider {
	return &HTTPProvider{}
}

// Name 提供者名称
func (h *HTTPProvider) Name() string {
	return "http"
}

// Register 注册服务
func (h *HTTPProvider) Register(application core.Application) error {
	// 获取配置和日志
	application.GetContainer().Invoke(func(config *conf.Config, logger log.Logger) {
		h.config = config
		h.logger = logger
	})

	if h.config == nil || h.logger == nil {
		return fmt.Errorf("无法获取配置或日志服务")
	}

	// 创建HTTP服务
	h.server = http.NewGinHTTPServer(h.config, h.logger)

	// 注册HTTP服务到容器
	application.GetContainer().Provide(func() http.HTTPServer {
		return h.server
	})

	// 注册Gin引擎到容器
	application.GetContainer().Provide(func() *gin.Engine {
		return h.server.Engine()
	})

	// 创建并设置路由适配器
	ginAdapter := &ginRouterAdapter{
		engine: h.server.Engine(),
	}

	// 设置全局路由
	facades.SetRoute(ginAdapter)

	h.logger.Info("HTTP服务已注册")
	return nil
}

// Boot 启动服务
func (h *HTTPProvider) Boot(application core.Application) error {
	// 无需在这里启动HTTP服务，由应用负责
	h.logger.Info("HTTP服务已初始化")
	return nil
}

// ginRouterAdapter Gin路由适配器
type ginRouterAdapter struct {
	engine *gin.Engine
}

// Group 创建路由组
func (g *ginRouterAdapter) Group(prefix string, handlers ...interface{}) facades.RouterGroup {
	// 直接返回Gin的RouterGroup
	return g.engine.Group(prefix)
}

// Run 启动服务
func (g *ginRouterAdapter) Run() error {
	return g.engine.Run()
}

// Shutdown 关闭服务
func (g *ginRouterAdapter) Shutdown() error {
	// Gin没有提供关闭方法，依赖HTTP服务关闭
	return nil
}
