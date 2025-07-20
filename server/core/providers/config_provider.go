package providers

import (
	"github.com/zhoudm1743/go-web/core"
	"github.com/zhoudm1743/go-web/core/conf"
)

// ConfigProvider 配置提供者
type ConfigProvider struct {
	config *conf.Config
}

// NewConfigProvider 创建配置提供者
func NewConfigProvider() *ConfigProvider {
	return &ConfigProvider{}
}

// Name 提供者名称
func (c *ConfigProvider) Name() string {
	return "config"
}

// Register 注册服务
func (c *ConfigProvider) Register(application core.Application) error {
	// 加载配置
	config, err := conf.NewConfig()
	if err != nil {
		return err
	}

	// 保存配置
	c.config = config

	// 注册到容器
	application.GetContainer().Provide(func() *conf.Config {
		return config
	})

	return nil
}

// Boot 启动服务
func (c *ConfigProvider) Boot(application core.Application) error {
	// 启动时不需要做任何事情
	return nil
}
