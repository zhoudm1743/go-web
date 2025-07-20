package controllers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-web/core/facades"
	"github.com/zhoudm1743/go-web/core/response"
	"github.com/zhoudm1743/go-web/pkg/generator"
)

// CodeGenController 代码生成器控制器
type CodeGenController struct{}

// NewCodeGenController 创建代码生成器控制器
func NewCodeGenController() *CodeGenController {
	return &CodeGenController{}
}

// GetApps 获取应用列表
func (c *CodeGenController) GetApps(ctx *gin.Context) {
	// 获取apps目录下的所有目录作为应用列表
	appsDir := "./server/apps"
	files, err := ioutil.ReadDir(appsDir)
	if err != nil {
		response.FailWithMsg(ctx, response.SystemError, fmt.Sprintf("获取应用列表失败: %v", err))
		return
	}

	// 过滤出目录
	var apps []string
	for _, file := range files {
		if file.IsDir() {
			apps = append(apps, file.Name())
		}
	}

	response.OkWithData(ctx, apps)
}

// GetTables 获取数据库表列表
func (c *CodeGenController) GetTables(ctx *gin.Context) {
	db := facades.DB()
	var tables []generator.TableInfo

	// 使用数据库原始查询
	// MySQL 查询
	query := `
		SELECT 
			table_name as tableName,
			table_comment as tableComment 
		FROM 
			information_schema.tables 
		WHERE 
			table_schema = DATABASE() 
		ORDER BY 
			table_name
	`

	if err := db.Raw(query).Scan(&tables).Error; err != nil {
		response.FailWithMsg(ctx, response.SystemError, fmt.Sprintf("获取表列表失败: %v", err))
		return
	}

	response.OkWithData(ctx, tables)
}

// GetColumns 获取表字段
func (c *CodeGenController) GetColumns(ctx *gin.Context) {
	tableName := ctx.Query("tableName")
	if tableName == "" {
		response.FailWithMsg(ctx, response.ParamsValidError, "表名不能为空")
		return
	}

	db := facades.DB()
	var columns []generator.ColumnInfo

	// MySQL 查询
	query := `
		SELECT 
			column_name as columnName,
			data_type as dataType,
			column_comment as columnComment,
			is_nullable as isNullable,
			column_key as columnKey
		FROM 
			information_schema.columns 
		WHERE 
			table_schema = DATABASE() 
			AND table_name = ? 
		ORDER BY 
			ordinal_position
	`

	if err := db.Raw(query, tableName).Scan(&columns).Error; err != nil {
		response.FailWithMsg(ctx, response.SystemError, fmt.Sprintf("获取表字段失败: %v", err))
		return
	}

	response.OkWithData(ctx, columns)
}

// Generate 生成代码
func (c *CodeGenController) Generate(ctx *gin.Context) {
	var req struct {
		StructName    string             `json:"structName"`
		TableName     string             `json:"tableName"`
		PackageName   string             `json:"packageName"`
		Description   string             `json:"description"`
		ApiPrefix     string             `json:"apiPrefix"`
		AppName       string             `json:"appName"`
		HasList       bool               `json:"hasList"`
		HasCreate     bool               `json:"hasCreate"`
		HasUpdate     bool               `json:"hasUpdate"`
		HasDelete     bool               `json:"hasDelete"`
		HasDetail     bool               `json:"hasDetail"`
		HasPagination bool               `json:"hasPagination"`
		Fields        []*generator.Field `json:"fields"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "请求参数有误")
		return
	}

	// 验证必填参数
	if req.StructName == "" {
		response.FailWithMsg(ctx, response.ParamsValidError, "结构体名称不能为空")
		return
	}

	if req.Description == "" {
		response.FailWithMsg(ctx, response.ParamsValidError, "描述不能为空")
		return
	}

	if len(req.Fields) == 0 {
		response.FailWithMsg(ctx, response.ParamsValidError, "至少需要一个字段")
		return
	}

	// 创建配置
	config := &generator.Config{
		StructName:    req.StructName,
		TableName:     req.TableName,
		PackageName:   req.PackageName,
		Description:   req.Description,
		ModuleName:    "github.com/zhoudm1743/go-web",
		RouterGroup:   "privateRoutes",
		ApiPrefix:     req.ApiPrefix,
		HasList:       req.HasList,
		HasCreate:     req.HasCreate,
		HasUpdate:     req.HasUpdate,
		HasDelete:     req.HasDelete,
		HasDetail:     req.HasDetail,
		HasPagination: req.HasPagination,
		Fields:        req.Fields,
	}

	// 设置项目根目录和应用目录
	rootPath := "./"
	appPath := filepath.Join("server/apps", req.AppName)

	// 检查应用目录是否存在，如果不存在则创建
	if _, err := os.Stat(appPath); os.IsNotExist(err) {
		// 创建应用目录结构
		dirs := []string{
			appPath,
			filepath.Join(appPath, "controllers"),
			filepath.Join(appPath, "models"),
			filepath.Join(appPath, "dto"),
			filepath.Join(appPath, "routes"),
			filepath.Join(appPath, "services"),
			filepath.Join(appPath, "middlewares"),
		}

		for _, dir := range dirs {
			if err := os.MkdirAll(dir, 0755); err != nil {
				response.FailWithMsg(ctx, response.SystemError, fmt.Sprintf("创建目录失败: %v", err))
				return
			}
		}

		// 创建基础文件
		appGoContent := fmt.Sprintf(`package %s

import (
	"github.com/gin-gonic/gin"
)

// Register 注册应用
func Register(r *gin.Engine) {
	// 初始化路由
	InitRoutes(r)
}
`, req.AppName)

		routesGoContent := fmt.Sprintf(`package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-web/apps/%s/controllers"
	"github.com/zhoudm1743/go-web/apps/%s/middlewares"
)

// InitRoutes 初始化路由
func InitRoutes(r *gin.Engine) {
	// 初始化控制器

	// 公开路由
	publicRoutes := r.Group("/%s")
	{
		// 路由示例
		// publicRoutes.GET("/example", exampleController.Example)
	}

	// 私有路由
	privateRoutes := r.Group("/%s")
	// privateRoutes.Use(middlewares.AuthMiddleware())
	{
		// 路由示例
		// privateRoutes.GET("/example", exampleController.Example)
	}
}
`, req.AppName, req.AppName, req.AppName, req.AppName)

		// 写入应用主文件
		appGoPath := filepath.Join(appPath, "app.go")
		if err := os.WriteFile(appGoPath, []byte(appGoContent), 0644); err != nil {
			response.FailWithMsg(ctx, response.SystemError, fmt.Sprintf("创建应用主文件失败: %v", err))
			return
		}

		// 写入路由文件
		routesGoPath := filepath.Join(appPath, "routes", "routes.go")
		if err := os.WriteFile(routesGoPath, []byte(routesGoContent), 0644); err != nil {
			response.FailWithMsg(ctx, response.SystemError, fmt.Sprintf("创建路由文件失败: %v", err))
			return
		}
	}

	// 创建生成器
	gen := generator.New(config)
	gen.SetRootPath(rootPath)

	// 初始化历史记录数据库
	if err := gen.InitHistoryDB(); err != nil {
		response.FailWithMsg(ctx, response.SystemError, fmt.Sprintf("初始化历史记录失败: %v", err))
		return
	}

	// 执行代码生成
	if err := gen.Run(); err != nil {
		response.FailWithMsg(ctx, response.SystemError, fmt.Sprintf("生成代码失败: %v", err))
		return
	}

	response.OkWithMsg(ctx, "生成代码成功")
}

// GetHistory 获取历史记录
func (c *CodeGenController) GetHistory(ctx *gin.Context) {
	// 解析分页参数
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	// 创建生成器
	gen := generator.New(&generator.Config{})
	gen.SetRootPath("./")

	// 获取历史记录列表
	list, total, err := gen.ListHistory(page, pageSize)
	if err != nil {
		response.FailWithMsg(ctx, response.SystemError, fmt.Sprintf("获取历史记录失败: %v", err))
		return
	}

	response.OkWithData(ctx, gin.H{
		"list":     list,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// Rollback 回滚代码生成
func (c *CodeGenController) Rollback(ctx *gin.Context) {
	var req struct {
		ID          uint `json:"id"`
		DeleteFiles bool `json:"deleteFiles"`
		DeleteAPI   bool `json:"deleteApi"`
		DeleteMenu  bool `json:"deleteMenu"`
		DeleteTable bool `json:"deleteTable"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "请求参数有误")
		return
	}

	if req.ID == 0 {
		response.FailWithMsg(ctx, response.ParamsValidError, "ID不能为空")
		return
	}

	// 创建生成器
	gen := generator.New(&generator.Config{})
	gen.SetRootPath("./")

	// 执行回滚
	if err := gen.RollBack(req.ID, req.DeleteFiles, req.DeleteAPI, req.DeleteMenu, req.DeleteTable); err != nil {
		response.FailWithMsg(ctx, response.SystemError, fmt.Sprintf("回滚失败: %v", err))
		return
	}

	response.OkWithMsg(ctx, "回滚成功")
}

// DeleteHistory 删除历史记录
func (c *CodeGenController) DeleteHistory(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "ID无效")
		return
	}

	// 创建生成器
	gen := generator.New(&generator.Config{})
	gen.SetRootPath("./")

	// 删除历史记录
	if err := gen.History.Delete(uint(id)); err != nil {
		response.FailWithMsg(ctx, response.SystemError, fmt.Sprintf("删除历史记录失败: %v", err))
		return
	}

	response.OkWithMsg(ctx, "删除成功")
}

// RegisterRoutes 注册路由
func (c *CodeGenController) RegisterRoutes(router *gin.RouterGroup) {
	codegenGroup := router.Group("/codegen")
	{
		codegenGroup.GET("/apps", c.GetApps)
		codegenGroup.GET("/tables", c.GetTables)
		codegenGroup.GET("/columns", c.GetColumns)
		codegenGroup.POST("/generate", c.Generate)
		codegenGroup.GET("/history", c.GetHistory)
		codegenGroup.POST("/rollback", c.Rollback)
		codegenGroup.DELETE("/history/:id", c.DeleteHistory)
	}
}
