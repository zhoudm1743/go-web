package cli

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-web/core/app"
)

// CLIApp 命令行应用
type CLIApp struct {
	*app.BaseApp
	commands []Command
}

// Command 命令接口
type Command interface {
	// Name 命令名称
	Name() string

	// Description 命令描述
	Description() string

	// Execute 执行命令
	Execute(args []string) error
}

// NewCLIApp 创建命令行应用
func NewCLIApp() *CLIApp {
	return &CLIApp{
		BaseApp:  app.NewBaseApp("cli"),
		commands: []Command{},
	}
}

// Initialize 初始化应用
func (a *CLIApp) Initialize() error {
	// 调用父类初始化
	if err := a.BaseApp.Initialize(); err != nil {
		return err
	}

	// 避免使用可能未初始化的facades
	// 使用fmt.Println代替facades.Log
	fmt.Println("初始化命令行应用:", a.Name())

	// 注册内置命令
	a.registerCommands()

	return nil
}

// RegisterRoutes 注册路由
// CLI应用不需要路由，但需要实现接口
func (a *CLIApp) RegisterRoutes(group *gin.RouterGroup) {
	// 命令行应用不需要路由
}

// Boot 启动应用
func (a *CLIApp) Boot() error {
	if err := a.BaseApp.Boot(); err != nil {
		return err
	}

	// 避免使用可能未初始化的facades
	fmt.Println("启动命令行应用:", a.Name())
	return nil
}

// Middlewares 获取应用中间件
func (a *CLIApp) Middlewares() []gin.HandlerFunc {
	// 命令行应用不需要中间件
	return []gin.HandlerFunc{}
}

// AddCommand 添加命令
func (a *CLIApp) AddCommand(cmd Command) {
	a.commands = append(a.commands, cmd)
}

// GetCommands 获取所有命令
func (a *CLIApp) GetCommands() []Command {
	return a.commands
}

// ExecuteCommand 执行命令
func (a *CLIApp) ExecuteCommand(name string, args []string) error {
	for _, cmd := range a.commands {
		if cmd.Name() == name {
			return cmd.Execute(args)
		}
	}

	// 避免使用可能未初始化的facades
	fmt.Printf("命令 %s 不存在\n", name)
	return nil
}

// registerCommands 注册内置命令
func (a *CLIApp) registerCommands() {
	// 这里可以注册内置命令
	// 例如：a.AddCommand(NewMigrateCommand())
}
