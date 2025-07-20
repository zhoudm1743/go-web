package providers

import (
	"github.com/zhoudm1743/go-web/core"
)

// AppProvider 应用提供者
type AppProvider struct{}

// NewAppProvider 创建应用提供者
func NewAppProvider() *AppProvider {
	return &AppProvider{}
}

// Name 提供者名称
func (p *AppProvider) Name() string {
	return "app"
}

// Register 注册服务到容器
func (p *AppProvider) Register(application core.Application) error {
	// AppProvider不再需要注册任何组件
	return nil
}

// Boot 启动服务
func (p *AppProvider) Boot(application core.Application) error {
	// AppProvider不再需要从容器中获取Manager
	return nil
}
