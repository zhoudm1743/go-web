package controllers

import (
	"fmt"
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
	// 获取工作目录
	workDir, err := os.Getwd()
	if err != nil {
		response.FailWithMsg(ctx, response.SystemError, fmt.Sprintf("获取工作目录失败: %v", err))
		return
	}

	// 检查可能的应用目录路径
	possiblePaths := []string{
		"apps",                          // 如果当前目录是server
		filepath.Join("server", "apps"), // 如果当前目录是项目根目录
	}

	var appsDir string
	var pathErr error

	for _, path := range possiblePaths {
		testPath := filepath.Join(workDir, path)
		if _, err := os.Stat(testPath); err == nil {
			appsDir = testPath
			pathErr = nil
			break
		} else {
			pathErr = err
		}
	}

	// 如果找不到有效路径，返回最后一个错误
	if appsDir == "" {
		response.FailWithMsg(ctx, response.SystemError, fmt.Sprintf("获取应用列表失败: %v", pathErr))
		return
	}

	files, err := os.ReadDir(appsDir)
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

	// 检测数据库类型，针对不同类型的数据库使用不同的查询语句
	dialect := db.Dialector.Name()
	fmt.Printf("当前数据库类型: %s\n", dialect)

	switch {
	case dialect == "sqlite" || dialect == "sqlite3":
		// 对于SQLite，直接查询表名，然后手动构建结果
		var tableNames []string
		query := `
			SELECT 
				name
			FROM 
				sqlite_master 
			WHERE 
				type = 'table' AND 
				name NOT LIKE 'sqlite_%' AND
				name NOT LIKE '%casbin%'
			ORDER BY 
				name
		`

		if err := db.Raw(query).Scan(&tableNames).Error; err != nil {
			response.FailWithMsg(ctx, response.SystemError, fmt.Sprintf("获取表列表失败: %v", err))
			return
		}

		// 手动构建结果数组
		tables := make([]generator.TableInfo, 0, len(tableNames))
		for _, name := range tableNames {
			tables = append(tables, generator.TableInfo{
				TableName:    name,
				TableComment: "",
			})
		}

		// 添加调试信息
		fmt.Printf("数据库类型: %s, 查询到 %d 个表:\n", dialect, len(tables))
		for i, table := range tables {
			fmt.Printf("  %d. 表名: %s\n", i+1, table.TableName)
		}

		response.OkWithData(ctx, tables)

	case dialect == "postgres" || dialect == "postgresql":
		// PostgreSQL查询
		var tables []generator.TableInfo
		query := `
			SELECT 
				table_name AS "tableName",
				obj_description((quote_ident(table_name)::text)::regclass, 'pg_class') AS "tableComment"
			FROM 
				information_schema.tables 
			WHERE 
				table_schema = 'public' AND 
				table_type = 'BASE TABLE' AND
				table_name NOT LIKE 'pg_%' AND
				table_name NOT LIKE '%casbin%'
			ORDER BY 
				table_name
		`

		if err := db.Raw(query).Scan(&tables).Error; err != nil {
			response.FailWithMsg(ctx, response.SystemError, fmt.Sprintf("获取表列表失败: %v", err))
			return
		}

		// 添加调试信息
		fmt.Printf("数据库类型: %s, 查询到 %d 个表:\n", dialect, len(tables))
		for i, table := range tables {
			fmt.Printf("  %d. 表名: %s, 注释: %s\n", i+1, table.TableName, table.TableComment)
		}

		response.OkWithData(ctx, tables)

	default: // MySQL和其他数据库
		// MySQL查询
		var tables []generator.TableInfo
		query := `
			SELECT 
				table_name AS "tableName",
				table_comment AS "tableComment" 
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

		// 添加调试信息
		fmt.Printf("数据库类型: %s, 查询到 %d 个表:\n", dialect, len(tables))
		for i, table := range tables {
			fmt.Printf("  %d. 表名: %s, 注释: %s\n", i+1, table.TableName, table.TableComment)
		}

		response.OkWithData(ctx, tables)
	}
}

// GetColumns 获取表字段
func (c *CodeGenController) GetColumns(ctx *gin.Context) {
	tableName := ctx.Query("tableName")
	if tableName == "" {
		response.FailWithMsg(ctx, response.ParamsValidError, "表名不能为空")
		return
	}

	db := facades.DB()

	// 检测数据库类型，针对不同类型的数据库使用不同的查询语句
	dialect := db.Dialector.Name()
	fmt.Printf("当前数据库类型: %s\n", dialect)

	switch {
	case dialect == "sqlite" || dialect == "sqlite3":
		// SQLite查询 - 使用pragma_table_info获取表结构
		type SqliteColumnInfo struct {
			Cid       int         `gorm:"column:cid"`
			Name      string      `gorm:"column:name"`
			Type      string      `gorm:"column:type"`
			NotNull   int         `gorm:"column:notnull"`
			DfltValue interface{} `gorm:"column:dflt_value"`
			Pk        int         `gorm:"column:pk"`
		}

		var sqliteColumns []SqliteColumnInfo
		query := `
			SELECT 
				cid, name, type, "notnull", dflt_value, pk
			FROM 
				pragma_table_info(?)
			ORDER BY 
				cid
		`

		if err := db.Raw(query, tableName).Scan(&sqliteColumns).Error; err != nil {
			response.FailWithMsg(ctx, response.SystemError, fmt.Sprintf("获取表字段失败: %v", err))
			return
		}

		// 转换为通用列信息格式
		columns := make([]generator.ColumnInfo, 0, len(sqliteColumns))
		for _, col := range sqliteColumns {
			isNullable := "YES"
			if col.NotNull == 1 {
				isNullable = "NO"
			}

			columnKey := ""
			if col.Pk == 1 {
				columnKey = "PRI"
			}

			columns = append(columns, generator.ColumnInfo{
				ColumnName:    col.Name,
				DataType:      col.Type,
				ColumnComment: "", // SQLite不支持列注释
				IsNullable:    isNullable,
				ColumnKey:     columnKey,
			})
		}

		// 添加调试信息
		fmt.Printf("表 %s 的字段信息 (%d 个):\n", tableName, len(columns))
		for i, col := range columns {
			fmt.Printf("  %d. 字段名: %s, 类型: %s, 主键: %s, 可空: %s\n",
				i+1, col.ColumnName, col.DataType, col.ColumnKey, col.IsNullable)
		}

		response.OkWithData(ctx, columns)

	case dialect == "postgres" || dialect == "postgresql":
		// PostgreSQL查询
		var columns []generator.ColumnInfo
		query := `
			SELECT 
				column_name AS "columnName",
				data_type AS "dataType",
				col_description(
					(table_schema || '.' || table_name)::regclass::oid, 
					ordinal_position
				) AS "columnComment",
				is_nullable AS "isNullable",
				CASE 
					WHEN EXISTS (
						SELECT 1 FROM information_schema.table_constraints tc
						JOIN information_schema.constraint_column_usage ccu 
						ON tc.constraint_name = ccu.constraint_name AND tc.table_schema = ccu.table_schema
						WHERE tc.constraint_type = 'PRIMARY KEY' 
						AND tc.table_name = ? 
						AND ccu.column_name = c.column_name
					) THEN 'PRI'
					ELSE ''
				END AS "columnKey"
			FROM 
				information_schema.columns c
			WHERE 
				table_schema = 'public' 
				AND table_name = ? 
			ORDER BY 
				ordinal_position
		`

		if err := db.Raw(query, tableName, tableName).Scan(&columns).Error; err != nil {
			response.FailWithMsg(ctx, response.SystemError, fmt.Sprintf("获取表字段失败: %v", err))
			return
		}

		// 添加调试信息
		fmt.Printf("表 %s 的字段信息 (%d 个):\n", tableName, len(columns))
		for i, col := range columns {
			fmt.Printf("  %d. 字段名: %s, 类型: %s, 主键: %s, 可空: %s\n",
				i+1, col.ColumnName, col.DataType, col.ColumnKey, col.IsNullable)
		}

		response.OkWithData(ctx, columns)

	default: // MySQL和其他数据库
		// MySQL查询
		var columns []generator.ColumnInfo
		query := `
			SELECT 
				column_name AS "columnName",
				data_type AS "dataType",
				column_comment AS "columnComment",
				is_nullable AS "isNullable",
				column_key AS "columnKey"
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

		// 添加调试信息
		fmt.Printf("表 %s 的字段信息 (%d 个):\n", tableName, len(columns))
		for i, col := range columns {
			fmt.Printf("  %d. 字段名: %s, 类型: %s, 主键: %s, 可空: %s\n",
				i+1, col.ColumnName, col.DataType, col.ColumnKey, col.IsNullable)
		}

		response.OkWithData(ctx, columns)
	}
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
	// 如果当前目录是server，则需要返回上一级
	if curDir, err := os.Getwd(); err == nil {
		if filepath.Base(curDir) == "server" {
			// 如果当前在server目录下，需要返回上一级
			rootPath = "../"
			fmt.Printf("当前目录是server，设置rootPath为: %s\n", rootPath)
		} else {
			fmt.Printf("当前工作目录: %s\n", curDir)
		}
	}
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
