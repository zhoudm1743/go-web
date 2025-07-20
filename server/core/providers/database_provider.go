package providers

import (
	"fmt"

	"github.com/zhoudm1743/go-web/core"
	"github.com/zhoudm1743/go-web/core/conf"
	"github.com/zhoudm1743/go-web/core/database"
	"github.com/zhoudm1743/go-web/core/facades"
	"github.com/zhoudm1743/go-web/core/log"
	"gorm.io/gorm"
)

// DatabaseProvider 数据库服务提供者
type DatabaseProvider struct{}

// NewDatabaseProvider 创建数据库服务提供者
func NewDatabaseProvider() *DatabaseProvider {
	return &DatabaseProvider{}
}

// Name 提供者名称
func (d *DatabaseProvider) Name() string {
	return "database"
}

// Register 注册服务到容器
func (d *DatabaseProvider) Register(app core.Application) error {
	container := app.GetContainer()

	// 先单独获取Config
	var config *conf.Config
	if err := container.Invoke(func(c *conf.Config) {
		config = c
	}); err != nil {
		return fmt.Errorf("无法获取配置: %w", err)
	}

	if config == nil {
		return fmt.Errorf("配置为空")
	}

	// 获取日志实例 - 如果不存在，创建一个默认的
	var logger log.Logger
	if err := container.Invoke(func(l log.Logger) {
		logger = l
	}); err != nil {
		// 如果无法获取日志实例，创建一个临时日志实例
		tempLogger, err := log.NewLogger(log.LoggerParams{Config: config})
		if err != nil {
			return fmt.Errorf("无法创建临时日志实例: %w", err)
		}
		logger = tempLogger
		fmt.Println("使用临时日志实例代替容器中的日志实例")
	}

	// 创建数据库连接
	db, err := database.NewDB(database.DBParams{
		Config: config,
		Logger: logger,
	})
	if err != nil {
		return fmt.Errorf("创建数据库连接失败: %w", err)
	}

	// 注册数据库到容器
	if err := container.Provide(func() *gorm.DB {
		return db
	}); err != nil {
		return fmt.Errorf("注册数据库实例到容器失败: %w", err)
	}

	// 同时注册GormDB结构体
	if err := container.Provide(func() *database.GormDB {
		return &database.GormDB{DB: db}
	}); err != nil {
		return fmt.Errorf("注册GormDB到容器失败: %w", err)
	}

	// 设置全局Facade
	facades.SetDB(db)

	fmt.Println("数据库提供者注册成功")
	return nil
}

// Boot 启动服务
func (d *DatabaseProvider) Boot(app core.Application) error {
	// 注册数据库关闭钩子
	// 注意：数据库可能未初始化，因此这里不进行任何操作
	return nil
}
