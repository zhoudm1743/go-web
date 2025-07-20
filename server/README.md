# Go-Web 框架

Go-Web 是一个简单高效的 Go Web 应用框架，结合了依赖注入(DI)和门面模式(Facades)的优点。这个框架旨在提供简洁直观的 API，同时保持高性能和可维护性。

## 特性

- **依赖注入**: 基于 `uber/dig` 实现的轻量级依赖注入系统
- **门面模式**: 提供优雅的全局访问点，便于使用各种服务
- **模块化**: 组件之间松耦合，易于扩展
- **高性能**: 基于 Gin 的 HTTP 路由
- **简单易用**: 简化的 API 和配置系统

## 安装

```bash
go get -u github.com/zhoudm1743/go-web
```

## 快速开始

### 创建 HTTP 服务

```go
package main

import (
	"github.com/zhoudm1743/go-web/core"
	"github.com/zhoudm1743/go-web/core/facades"
	"github.com/zhoudm1743/go-web/core/providers"
)

func main() {
	// 创建应用实例
	app := core.NewApp()
	facades.SetApp(app)

	// 注册服务提供者
	app.Register(providers.NewConfigProvider())
	app.Register(providers.NewLogProvider())
	app.Register(providers.NewHTTPProvider())
	
	// 启动应用
	app.Boot()

	// 定义路由
	facades.Route().Group("/api").GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	// 启动HTTP服务
	go facades.Route().Run()
	
	// 运行应用（阻塞，直到收到退出信号）
	app.Run()
}
```

### 使用数据库

```go
// 注册数据库服务提供者
app.Register(providers.NewDatabaseProvider())

// 使用门面模式访问数据库
db := facades.DB()

// 创建用户
user := &models.User{Name: "张三"}
db.Create(user)
```

### 使用缓存

```go
// 注册缓存服务提供者
app.Register(providers.NewCacheProvider())

// 使用门面模式访问缓存
facades.Cache().Set("key", "value", 10*time.Minute)

// 获取缓存
value, err := facades.Cache().Get("key")
```

## 核心概念

### 应用生命周期

1. 创建应用实例 (`core.NewApp()`)
2. 注册服务提供者 (`app.Register()`)
3. 启动应用 (`app.Boot()`)
4. 运行应用 (`app.Run()`)

### 服务提供者

服务提供者负责绑定服务到容器中，并在应用启动时进行初始化。每个服务提供者需要实现以下接口：

```go
type Provider interface {
	// 提供者名称
	Name() string
	// 注册服务到容器
	Register(app Application) error
	// 启动服务
	Boot(app Application) error
}
```

### 门面模式

框架提供了多种门面，用于在全局范围内访问服务：

- `facades.Log()`: 日志服务
- `facades.DB()`: 数据库服务
- `facades.Cache()`: 缓存服务
- `facades.Route()`: HTTP路由服务
- `facades.Config()`: 配置服务

## 配置

配置文件位于 `config` 目录，采用 YAML 格式。主要配置项包括：

- App: 应用基本配置
- HTTP: HTTP服务配置
- Database: 数据库配置
- Log: 日志配置
- Cache: 缓存配置

## 扩展

### 创建自定义服务提供者

```go
package providers

import (
	"github.com/zhoudm1743/go-web/core"
)

// MyServiceProvider 自定义服务提供者
type MyServiceProvider struct{}

func NewMyServiceProvider() *MyServiceProvider {
	return &MyServiceProvider{}
}

func (m *MyServiceProvider) Name() string {
	return "my-service"
}

func (m *MyServiceProvider) Register(app core.Application) error {
	container := app.GetContainer()
	
	// 注册服务到容器
	return container.Provide(func() *MyService {
		return NewMyService()
	})
}

func (m *MyServiceProvider) Boot(app core.Application) error {
	// 服务启动逻辑
	return nil
}
```

## 许可证

MIT 