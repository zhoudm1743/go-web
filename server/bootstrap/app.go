package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-web/apps/admin"
	"github.com/zhoudm1743/go-web/apps/api"
	"github.com/zhoudm1743/go-web/apps/cli"
	"github.com/zhoudm1743/go-web/core"
	"github.com/zhoudm1743/go-web/core/app"
	"github.com/zhoudm1743/go-web/core/conf"
	"github.com/zhoudm1743/go-web/core/database"
	"github.com/zhoudm1743/go-web/core/facades"
	"github.com/zhoudm1743/go-web/core/log"
	"github.com/zhoudm1743/go-web/core/utils"
	"github.com/zhoudm1743/go-web/routes"
)

// Application 应用实例
type Application struct {
	core    *core.App          // 核心应用实例
	config  *conf.Config       // 配置
	logger  log.Logger         // 日志
	engine  *gin.Engine        // Gin引擎
	server  *http.Server       // HTTP服务器
	manager *app.Manager       // 应用管理器
	appMode string             // 应用模式：http, cli
	signal  chan os.Signal     // 信号通道
	ctx     context.Context    // 上下文
	cancel  context.CancelFunc // 取消函数
}

// NewApplication 创建应用实例
func NewApplication() (*Application, error) {
	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())

	// 创建应用
	app := &Application{
		core:    core.NewApp(),
		signal:  make(chan os.Signal, 1),
		ctx:     ctx,
		cancel:  cancel,
		appMode: "http", // 默认为HTTP模式
	}

	// 注册信号处理
	signal.Notify(app.signal, syscall.SIGINT, syscall.SIGTERM)

	return app, nil
}

// Initialize 初始化应用
func (a *Application) Initialize() error {
	// 加载配置
	config, err := conf.NewConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}
	a.config = config

	// 设置全局配置
	facades.SetConfig(config)

	// 设置Gin模式
	switch a.config.App.Mode {
	case "prod", "production":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	// 注册服务提供者
	if err := registerProviders(a.core); err != nil {
		return fmt.Errorf("注册服务提供者失败: %w", err)
	}

	// 创建一个临时的控制台日志器，避免在初始化过程中出现nil指针
	tempLogger, err := log.NewLogger(log.LoggerParams{Config: a.config})
	if err == nil {
		a.logger = tempLogger
	} else {
		fmt.Printf("警告: 无法创建临时日志器: %v\n", err)
	}

	// 尝试从全局facades获取日志实例
	if facades.Log() != nil {
		a.logger = facades.Log()
		fmt.Println("成功从facades获取日志实例")
	}

	// 创建应用管理器
	a.manager = app.NewManager(a.core.GetContainer())

	// 注册应用
	if err := a.registerApplications(); err != nil {
		return fmt.Errorf("注册应用失败: %w", err)
	}

	// 设置到全局Facades
	facades.SetAppManager(a.manager)

	// 初始化数据库相关内容
	if err := a.initDatabase(); err != nil {
		return fmt.Errorf("初始化数据库失败: %w", err)
	}

	// 初始化所有应用
	if err := a.manager.InitializeApps(); err != nil {
		return fmt.Errorf("初始化应用失败: %w", err)
	}

	// 创建Gin引擎
	a.engine = a.createGinEngine()

	// 注册路由
	a.registerRoutes()

	// 创建HTTP服务器
	a.server = &http.Server{
		Addr:           fmt.Sprintf("%s:%d", a.config.HTTP.Host, a.config.HTTP.Port),
		Handler:        a.engine,
		ReadTimeout:    a.config.HTTP.ReadTimeout,
		WriteTimeout:   a.config.HTTP.WriteTimeout,
		MaxHeaderBytes: a.config.HTTP.MaxHeaderBytes,
	}

	return nil
}

// initDatabase 初始化数据库相关内容
func (a *Application) initDatabase() error {
	// 先尝试从facades获取数据库实例
	db := facades.DB()

	if db == nil {
		// 使用database provider直接创建一个新的数据库实例
		dbProvider := a.core.GetProvider("database")
		if dbProvider == nil {
			return fmt.Errorf("找不到数据库提供者")
		}

		// 手动创建数据库连接
		dbInstance, err := database.NewDB(database.DBParams{
			Config: a.config,
			Logger: a.logger,
		})
		if err != nil {
			return fmt.Errorf("手动创建数据库连接失败: %w", err)
		}

		// 设置到facades
		facades.SetDB(dbInstance)
		db = dbInstance
	}

	// 初始化Casbin表和策略
	if err := utils.InitCasbinTables(db); err != nil {
		return fmt.Errorf("初始化Casbin表和策略失败: %w", err)
	}

	return nil
}

// registerApplications 注册应用
func (a *Application) registerApplications() error {
	// 注册管理后台应用
	if err := a.manager.RegisterApp(admin.NewApp()); err != nil {
		return err
	}

	// 注册API应用
	if err := a.manager.RegisterApp(api.NewAPIApp()); err != nil {
		return err
	}

	// 注册命令行应用
	if err := a.manager.RegisterApp(cli.NewCLIApp()); err != nil {
		return err
	}

	return nil
}

// createGinEngine 创建Gin引擎
func (a *Application) createGinEngine() *gin.Engine {
	engine := gin.New()

	// 添加恢复中间件
	engine.Use(gin.Recovery())

	// 设置日志中间件
	if a.logger != nil {
		// 可以根据需要自定义日志中间件
		engine.Use(func(c *gin.Context) {
			start := time.Now()
			path := c.Request.URL.Path
			raw := c.Request.URL.RawQuery

			c.Next()

			latency := time.Since(start)
			statusCode := c.Writer.Status()
			clientIP := c.ClientIP()
			method := c.Request.Method

			a.logger.Infof("%s %s %d %s %s %s",
				clientIP, method, statusCode, path, raw, latency)
		})
	}

	return engine
}

// registerRoutes 注册路由
func (a *Application) registerRoutes() {
	// 注册全局中间件
	routes.RegisterGlobalMiddlewares(a.engine)

	// 注册全局路由
	routes.RegisterGlobalRoutes(a.engine)

	// 为每个应用注册路由
	for _, appInstance := range a.manager.GetApps() {
		appName := appInstance.Name()

		// 创建应用路由组
		group := a.engine.Group("/" + appName)

		// 注册应用中间件
		for _, middleware := range appInstance.Middlewares() {
			group.Use(middleware)
		}

		// 注册应用路由
		appInstance.RegisterRoutes(group)

		a.logger.Infof("应用 [%s] 路由已注册", appName)
	}
}

// SetMode 设置应用模式
func (a *Application) SetMode(mode string) {
	a.appMode = mode
}

// Run 运行应用
func (a *Application) Run() error {
	// 启动所有应用
	if err := a.manager.BootApps(); err != nil {
		return fmt.Errorf("启动应用失败: %w", err)
	}

	// 根据模式运行
	switch a.appMode {
	case "cli":
		// CLI模式下，可能没有需要运行的HTTP服务
		a.waitForSignal()
	default:
		// HTTP模式
		go func() {
			a.logger.Infof("HTTP服务已启动: %s:%d", a.config.HTTP.Host, a.config.HTTP.Port)
			if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				a.logger.Errorf("HTTP服务启动失败: %v", err)
			}
		}()

		// 等待信号
		a.waitForSignal()
	}

	return nil
}

// waitForSignal 等待信号
func (a *Application) waitForSignal() {
	// 等待中断信号
	<-a.signal
	a.logger.Info("收到关闭信号，正在关闭应用...")

	// 执行关闭逻辑
	a.Shutdown()
}

// Shutdown 关闭应用
func (a *Application) Shutdown() error {
	// 取消上下文
	a.cancel()

	// 创建超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 关闭HTTP服务
	if a.server != nil {
		if err := a.server.Shutdown(ctx); err != nil {
			a.logger.Errorf("HTTP服务关闭失败: %v", err)
		} else {
			a.logger.Info("HTTP服务已关闭")
		}
	}

	// 调用应用实例的Shutdown方法
	if err := a.core.Shutdown(); err != nil {
		a.logger.Errorf("应用关闭失败: %v", err)
	}

	a.logger.Info("应用已完全关闭")
	return nil
}

// GetCore 获取核心应用实例
func (a *Application) GetCore() *core.App {
	return a.core
}

// GetEngine 获取Gin引擎
func (a *Application) GetEngine() *gin.Engine {
	return a.engine
}

// GetManager 获取应用管理器
func (a *Application) GetManager() *app.Manager {
	return a.manager
}

// InitializeApp 初始化应用
func InitializeApp() (*Application, error) {
	app, err := NewApplication()
	if err != nil {
		return nil, err
	}

	if err := app.Initialize(); err != nil {
		return nil, err
	}

	return app, nil
}
