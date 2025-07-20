package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-web/core/conf"
	"github.com/zhoudm1743/go-web/core/log"
)

// HTTPServer HTTP服务接口
type HTTPServer interface {
	// Run 启动HTTP服务
	Run() error
	// Shutdown 关闭HTTP服务
	Shutdown(ctx context.Context) error
	// Engine 返回底层Gin引擎
	Engine() *gin.Engine
}

// GinHTTPServer 基于Gin的HTTP服务实现
type GinHTTPServer struct {
	server *http.Server
	engine *gin.Engine
	config *conf.HTTPConfig
	Logger log.Logger
}

// NewGinHTTPServer 创建Gin HTTP服务
func NewGinHTTPServer(config *conf.Config, logger log.Logger) *GinHTTPServer {
	// 创建自定义Gin引擎，使用我们的日志系统
	engine := log.CustomGinEngine(logger, config.App.Mode == "prod")

	// 配置HTTP服务器
	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", config.HTTP.Host, config.HTTP.Port),
		Handler:        engine,
		ReadTimeout:    config.HTTP.ReadTimeout,
		WriteTimeout:   config.HTTP.WriteTimeout,
		MaxHeaderBytes: config.HTTP.MaxHeaderBytes,
	}

	return &GinHTTPServer{
		server: server,
		engine: engine,
		config: &config.HTTP,
		Logger: logger,
	}
}

// Engine 返回底层Gin引擎
func (g *GinHTTPServer) Engine() *gin.Engine {
	return g.engine
}

// Run 启动HTTP服务
func (g *GinHTTPServer) Run() error {
	g.Logger.Info("HTTP服务和路由已初始化")
	// 启动HTTP服务
	return g.server.ListenAndServe()
}

// Shutdown 关闭HTTP服务
func (g *GinHTTPServer) Shutdown(ctx context.Context) error {
	// 关闭HTTP服务
	return g.server.Shutdown(ctx)
}
