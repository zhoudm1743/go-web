package providers

import (
	"fmt"

	"github.com/zhoudm1743/go-web/core"
	"github.com/zhoudm1743/go-web/core/conf"
	"github.com/zhoudm1743/go-web/core/facades"
	"github.com/zhoudm1743/go-web/core/log"
)

// LogProvider 日志服务提供者
type LogProvider struct {
	logger log.Logger
}

// NewLogProvider 创建日志服务提供者
func NewLogProvider() *LogProvider {
	return &LogProvider{}
}

// Name 提供者名称
func (l *LogProvider) Name() string {
	return "log"
}

// Register 注册服务到容器
func (l *LogProvider) Register(app core.Application) error {
	// 直接从容器获取Config
	var config *conf.Config
	if err := app.GetContainer().Invoke(func(c *conf.Config) {
		config = c
	}); err != nil {
		return fmt.Errorf("无法获取配置: %w", err)
	}

	if config == nil {
		return fmt.Errorf("配置为空")
	}

	// 创建日志实例
	logger, err := log.NewLogger(log.LoggerParams{
		Config: config,
	})
	if err != nil {
		return fmt.Errorf("创建日志实例失败: %w", err)
	}

	// 保存实例
	l.logger = logger

	// 使用正确的方式注册到容器
	if err := app.GetContainer().Provide(func() log.Logger {
		return logger
	}); err != nil {
		return fmt.Errorf("注册日志实例失败: %w", err)
	}

	// 设置全局Facade
	facades.SetLog(logger)

	fmt.Println("日志服务已注册")
	return nil
}

// Boot 启动服务
func (l *LogProvider) Boot(app core.Application) error {
	if l.logger != nil {
		l.logger.Info("日志服务已启动")
	} else {
		fmt.Println("日志服务已启动（日志实例为空）")
	}
	return nil
}
