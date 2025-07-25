package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-web/apps/admin/models"
	"github.com/zhoudm1743/go-web/apps/admin/routes"
	"github.com/zhoudm1743/go-web/core/app"
	"github.com/zhoudm1743/go-web/core/facades"
	"github.com/zhoudm1743/go-web/core/utils"
)

// App 管理后台应用
type App struct {
	app.BaseApp
}

// NewApp 创建管理后台应用
func NewApp() *App {
	return &App{}
}

// Name 应用名称
func (a *App) Name() string {
	return "admin"
}

// Initialize 初始化应用
func (a *App) Initialize() error {
	// 自动迁移数据库表结构
	db := facades.DB()
	if db == nil {
		return nil
	}

	// 自动迁移模型到数据库表
	err := db.AutoMigrate(
		&models.Admin{},
		&models.Role{},
		&models.Menu{},
		&models.RoleMenu{},
		&models.AdminRole{},
	)
	if err != nil {
		return err
	}

	// 初始化默认管理员账号
	if err := models.CreateDefaultAdminIfNotExists(db); err != nil {
		return err
	}

	// 初始化Casbin表和权限
	if err := utils.InitCasbinTables(db); err != nil {
		return err
	}

	return nil
}

// RegisterRoutes 注册路由
func (a *App) RegisterRoutes(r *gin.RouterGroup) {
	// 注册应用的所有路由
	routes.RegisterRoutes(r)
}

// Boot 启动应用
func (a *App) Boot() error {
	return nil
}

// Middlewares 中间件
func (a *App) Middlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{}
}
