package generator

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-web/core/response"
)

// CodeGenController 代码生成器控制器
type CodeGenController struct {
	rootPath string // 项目根目录
}

// NewCodeGenController 创建代码生成器控制器
func NewCodeGenController(rootPath string) *CodeGenController {
	return &CodeGenController{
		rootPath: rootPath,
	}
}

// GenerateCode 生成代码
func (c *CodeGenController) GenerateCode(ctx *gin.Context) {
	// 解析请求参数
	var req struct {
		StructName    string   `json:"structName"`
		TableName     string   `json:"tableName"`
		PackageName   string   `json:"packageName"`
		Description   string   `json:"description"`
		ApiPrefix     string   `json:"apiPrefix"`
		HasList       bool     `json:"hasList"`
		HasCreate     bool     `json:"hasCreate"`
		HasUpdate     bool     `json:"hasUpdate"`
		HasDelete     bool     `json:"hasDelete"`
		HasDetail     bool     `json:"hasDetail"`
		HasPagination bool     `json:"hasPagination"`
		Fields        []*Field `json:"fields"`
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
	config := &Config{
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

	// 设置默认值
	if config.PackageName == "" {
		config.PackageName = "admin"
	}
	if config.TableName == "" {
		config.TableName = ToSnakeCase(config.StructName)
	}
	if config.ApiPrefix == "" {
		config.ApiPrefix = config.TableName
	}

	// 创建生成器
	gen := New(config)
	gen.SetRootPath(c.rootPath)

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

// GetHistoryList 获取历史记录列表
func (c *CodeGenController) GetHistoryList(ctx *gin.Context) {
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
	gen := New(&Config{})
	gen.SetRootPath(c.rootPath)

	// 获取历史记录列表
	list, total, err := gen.ListHistory(page, pageSize)
	if err != nil {
		response.FailWithMsg(ctx, response.SystemError, fmt.Sprintf("获取历史记录列表失败: %v", err))
		return
	}

	response.OkWithData(ctx, gin.H{
		"list":     list,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// RollBack 回滚代码生成
func (c *CodeGenController) RollBack(ctx *gin.Context) {
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
	gen := New(&Config{})
	gen.SetRootPath(c.rootPath)

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
	gen := New(&Config{})
	gen.SetRootPath(c.rootPath)

	// 删除历史记录
	if err := gen.History.Delete(uint(id)); err != nil {
		response.FailWithMsg(ctx, response.SystemError, fmt.Sprintf("删除历史记录失败: %v", err))
		return
	}

	response.OkWithMsg(ctx, "删除成功")
}

// GetDBTables 获取数据库表列表
func (c *CodeGenController) GetDBTables(ctx *gin.Context) {
	// 这里需要根据您的实际情况实现获取数据库表的逻辑
	// 为简化示例，这里返回模拟数据
	tables := []map[string]interface{}{
		{"tableName": "users", "tableComment": "用户表"},
		{"tableName": "roles", "tableComment": "角色表"},
		{"tableName": "permissions", "tableComment": "权限表"},
	}

	response.OkWithData(ctx, gin.H{
		"tables": tables,
	})
}

// GetTableColumns 获取表字段
func (c *CodeGenController) GetTableColumns(ctx *gin.Context) {
	tableName := ctx.Query("tableName")
	if tableName == "" {
		response.FailWithMsg(ctx, response.ParamsValidError, "表名不能为空")
		return
	}

	// 这里需要根据您的实际情况实现获取表字段的逻辑
	// 为简化示例，这里返回模拟数据
	columns := []map[string]interface{}{
		{
			"columnName":    "id",
			"dataType":      "int",
			"columnComment": "主键ID",
			"isNullable":    "NO",
			"columnKey":     "PRI",
		},
		{
			"columnName":    "name",
			"dataType":      "varchar",
			"columnComment": "名称",
			"isNullable":    "NO",
			"columnKey":     "",
		},
		{
			"columnName":    "status",
			"dataType":      "tinyint",
			"columnComment": "状态",
			"isNullable":    "YES",
			"columnKey":     "",
		},
	}

	response.OkWithData(ctx, gin.H{
		"columns": columns,
	})
}

// RegisterRoutes 注册路由
func (c *CodeGenController) RegisterRoutes(router *gin.RouterGroup) {
	codegenGroup := router.Group("/codegen")
	{
		codegenGroup.POST("/generate", c.GenerateCode)
		codegenGroup.GET("/history", c.GetHistoryList)
		codegenGroup.POST("/rollback", c.RollBack)
		codegenGroup.DELETE("/history/:id", c.DeleteHistory)
		codegenGroup.GET("/tables", c.GetDBTables)
		codegenGroup.GET("/columns", c.GetTableColumns)
	}
}

// ToSnakeCase 将驼峰命名转换为蛇形命名
func ToSnakeCase(s string) string {
	return ToLowerCamel(s)
}
