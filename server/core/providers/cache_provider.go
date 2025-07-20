package providers

import (
	"github.com/zhoudm1743/go-web/core"
	"github.com/zhoudm1743/go-web/core/cache"
	"github.com/zhoudm1743/go-web/core/conf"
	"github.com/zhoudm1743/go-web/core/facades"
	"github.com/zhoudm1743/go-web/core/log"
)

// CacheProvider 缓存服务提供者
type CacheProvider struct{}

// NewCacheProvider 创建缓存服务提供者
func NewCacheProvider() *CacheProvider {
	return &CacheProvider{}
}

// Name 提供者名称
func (c *CacheProvider) Name() string {
	return "cache"
}

// Register 注册服务到容器
func (c *CacheProvider) Register(app core.Application) error {
	container := app.GetContainer()

	// 使用DI解析依赖
	return container.Invoke(func(config *conf.Config, logger log.Logger) error {
		// 根据配置创建缓存实例
		var cacheInstance cache.Cache
		var err error

		switch config.Cache.Type {
		case "redis":
			// 创建Redis缓存
			cacheInstance, err = cache.NewRedisCache(config, logger)
		case "file":
			// 创建文件缓存
			cacheInstance, err = cache.NewFileCache(config, logger)
		default:
			// 默认使用内存缓存
			cacheInstance, err = cache.NewMemoryCache(config, logger)
		}

		if err != nil {
			return err
		}

		// 注册缓存到容器
		if err := container.Provide(func() cache.Cache {
			return cacheInstance
		}); err != nil {
			return err
		}

		// 设置全局Facade
		facades.SetCache(cacheInstance)

		return nil
	})
}

// Boot 启动服务
func (c *CacheProvider) Boot(app core.Application) error {
	// 缓存服务不需要启动逻辑
	return nil
}
