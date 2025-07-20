// Package app provides application management functionality
package app

import (
	"fmt"
	"sync"

	"go.uber.org/dig"
)

// Manager 应用管理器实现
type Manager struct {
	apps      map[string]AppInterface
	container *dig.Container
	mu        sync.RWMutex
}

// 确保Manager实现AppManager接口
var _ AppManager = (*Manager)(nil)

// NewManager 创建应用管理器
func NewManager(container *dig.Container) *Manager {
	return &Manager{
		apps:      make(map[string]AppInterface),
		container: container,
		mu:        sync.RWMutex{},
	}
}

// RegisterApp 注册应用
func (m *Manager) RegisterApp(app AppInterface) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	name := app.Name()
	if _, exists := m.apps[name]; exists {
		return fmt.Errorf("应用 %s 已存在", name)
	}

	m.apps[name] = app
	return nil
}

// GetApp 获取应用
func (m *Manager) GetApp(name string) (AppInterface, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	app, ok := m.apps[name]
	return app, ok
}

// GetApps 获取所有应用
func (m *Manager) GetApps() []AppInterface {
	m.mu.RLock()
	defer m.mu.RUnlock()

	apps := make([]AppInterface, 0, len(m.apps))
	for _, app := range m.apps {
		apps = append(apps, app)
	}
	return apps
}

// InitializeApps 初始化所有应用
func (m *Manager) InitializeApps() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for name, app := range m.apps {
		if err := app.Initialize(); err != nil {
			return fmt.Errorf("初始化应用 %s 失败: %w", name, err)
		}
	}
	return nil
}

// BootApps 启动所有应用
func (m *Manager) BootApps() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for name, app := range m.apps {
		if err := app.Boot(); err != nil {
			return fmt.Errorf("启动应用 %s 失败: %w", name, err)
		}
	}
	return nil
}

// ShutdownApps 关闭所有应用
func (m *Manager) ShutdownApps() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var lastErr error
	for name, app := range m.apps {
		if err := app.Shutdown(); err != nil {
			lastErr = fmt.Errorf("关闭应用 %s 失败: %w", name, err)
		}
	}
	return lastErr
}
