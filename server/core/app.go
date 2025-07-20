package core

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/dig"
)

// Application 应用框架接口
type Application interface {
	// Boot 启动应用
	Boot() error
	// Register 注册组件
	Register(provider Provider) error
	// GetContainer 获取DI容器
	GetContainer() *dig.Container
	// GetProviders 获取所有提供者
	GetProviders() []Provider
	// GetProvider 通过名称获取提供者
	GetProvider(name string) Provider
	// Run 运行应用
	Run() error
	// Shutdown 优雅关闭
	Shutdown() error
}

// Provider 服务提供者接口
type Provider interface {
	// Name 提供者名称
	Name() string
	// Register 注册服务到容器
	Register(app Application) error
	// Boot 启动服务
	Boot(app Application) error
}

// App 应用实现
type App struct {
	container *dig.Container
	providers []Provider
	booted    bool
}

// NewApp 创建应用实例
func NewApp() *App {
	return &App{
		container: dig.New(),
		providers: []Provider{},
		booted:    false,
	}
}

// Boot 启动应用
func (a *App) Boot() error {
	if a.booted {
		return nil
	}

	// 注册所有提供者
	for _, provider := range a.providers {
		if err := provider.Register(a); err != nil {
			return err
		}
	}

	// 启动所有提供者
	for _, provider := range a.providers {
		if err := provider.Boot(a); err != nil {
			return err
		}
	}

	a.booted = true
	return nil
}

// Register 注册组件
func (a *App) Register(provider Provider) error {
	a.providers = append(a.providers, provider)
	return nil
}

// GetContainer 获取DI容器
func (a *App) GetContainer() *dig.Container {
	return a.container
}

// GetProviders 获取所有提供者
func (a *App) GetProviders() []Provider {
	return a.providers
}

// GetProvider 通过名称获取提供者
func (a *App) GetProvider(name string) Provider {
	for _, provider := range a.providers {
		if provider.Name() == name {
			return provider
		}
	}
	return nil
}

// Run 运行应用
func (a *App) Run() error {
	// 如果应用未启动，先启动
	if !a.booted {
		if err := a.Boot(); err != nil {
			return err
		}
	}

	// 创建信号监听通道
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 等待信号
	<-quit

	// 关闭应用
	return a.Shutdown()
}

// Shutdown 优雅关闭
func (a *App) Shutdown() error {
	// 这里会调用各组件的关闭逻辑
	// 目前是一个简单实现，后续可以添加更多的关闭逻辑
	return nil
}
