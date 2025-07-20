package bootstrap

import (
	"fmt"

	"github.com/zhoudm1743/go-web/core"
	"github.com/zhoudm1743/go-web/core/providers"
)

// registerProviders 注册服务提供者
func registerProviders(app core.Application) error {
	// 注册配置提供者（必须在其他提供者之前注册）
	if err := app.Register(providers.NewConfigProvider()); err != nil {
		return fmt.Errorf("注册配置提供者失败: %w", err)
	}

	// 立即启动配置提供者
	configProvider := app.GetProvider("config")
	if configProvider != nil {
		if err := configProvider.Boot(app); err != nil {
			return fmt.Errorf("启动配置提供者失败: %w", err)
		}
		fmt.Println("配置提供者已启动")
	} else {
		return fmt.Errorf("配置提供者未注册")
	}

	// 注册日志提供者
	if err := app.Register(providers.NewLogProvider()); err != nil {
		return fmt.Errorf("注册日志提供者失败: %w", err)
	}

	// 立即启动日志提供者
	logProvider := app.GetProvider("log")
	if logProvider != nil {
		if err := logProvider.Boot(app); err != nil {
			return fmt.Errorf("启动日志提供者失败: %w", err)
		}
		fmt.Println("日志提供者已启动")
	} else {
		return fmt.Errorf("日志提供者未注册")
	}

	// 注册数据库提供者
	if err := app.Register(providers.NewDatabaseProvider()); err != nil {
		return fmt.Errorf("注册数据库提供者失败: %w", err)
	}

	// 注册缓存提供者
	if err := app.Register(providers.NewCacheProvider()); err != nil {
		return fmt.Errorf("注册缓存提供者失败: %w", err)
	}

	// 启动数据库和缓存提供者
	for _, name := range []string{"database", "cache"} {
		provider := app.GetProvider(name)
		if provider != nil {
			if err := provider.Boot(app); err != nil {
				return fmt.Errorf("启动%s提供者失败: %w", name, err)
			}
			fmt.Printf("%s提供者已启动\n", name)
		} else {
			return fmt.Errorf("%s提供者未注册", name)
		}
	}

	// 注册应用提供者
	if err := app.Register(providers.NewAppProvider()); err != nil {
		return fmt.Errorf("注册应用提供者失败: %w", err)
	}

	// 启动应用提供者
	appProvider := app.GetProvider("app")
	if appProvider != nil {
		if err := appProvider.Boot(app); err != nil {
			return fmt.Errorf("启动应用提供者失败: %w", err)
		}
		fmt.Println("应用提供者已启动")
	} else {
		return fmt.Errorf("应用提供者未注册")
	}

	// 注册HTTP服务提供者 - 移到最后，确保所有其他服务都已注册
	if err := app.Register(providers.NewHTTPProvider()); err != nil {
		return fmt.Errorf("注册HTTP服务提供者失败: %w", err)
	}

	// 在这里注册更多服务提供者

	// 检查各提供者状态
	for _, provider := range app.GetProviders() {
		fmt.Printf("提供者 %s 已注册\n", provider.Name())
	}

	return nil
}
