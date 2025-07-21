package generator

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"
)

// // Config 代码生成器配置
// type Config struct {
// 	// 基本信息
// 	StructName  string // 结构体名称
// 	TableName   string // 表名
// 	PackageName string // 包名
// 	Description string // 描述
// 	ModuleName  string // 模块名(import路径前缀)

// 	// 路由和API
// 	RouterGroup string // 路由分组
// 	ApiPrefix   string // API前缀

// 	// 选项
// 	HasList       bool // 是否有列表
// 	HasCreate     bool // 是否有创建
// 	HasUpdate     bool // 是否有更新
// 	HasDelete     bool // 是否有删除
// 	HasDetail     bool // 是否有详情
// 	HasPagination bool // 是否分页

// 	// 字段配置
// 	Fields []*Field
// }

// // Field 字段定义
// type Field struct {
// 	FieldName    string // 结构体字段名称
// 	FieldType    string // Go类型
// 	ColumnName   string // 数据库字段名
// 	ColumnType   string // 数据库字段类型
// 	FieldDesc    string // 字段描述
// 	Required     bool   // 是否必填
// 	IsPrimaryKey bool   // 是否主键
// 	IsSearchable bool   // 是否可搜索
// 	IsFilterable bool   // 是否可过滤
// 	IsSortable   bool   // 是否可排序
// }

// Generator 代码生成器
type Generator struct {
	Config         *Config
	RootPath       string            // 项目根目录
	History        *HistoryManager   // 历史记录管理器
	generatedFiles map[string]string // 生成的文件路径映射
}

// New 创建代码生成器
func New(config *Config) *Generator {
	return &Generator{
		Config:         config,
		RootPath:       "./",
		History:        NewHistoryManager(),
		generatedFiles: make(map[string]string),
	}
}

// SetRootPath 设置项目根目录
func (g *Generator) SetRootPath(path string) {
	g.RootPath = path
}

// InitHistoryDB 初始化历史记录数据库
func (g *Generator) InitHistoryDB() error {
	return g.History.Migrate()
}

// Run 执行代码生成
func (g *Generator) Run() error {
	// 生成模型
	if err := g.generateModel(); err != nil {
		return err
	}

	// 生成DTO
	if err := g.generateDTO(); err != nil {
		return err
	}

	// 生成控制器
	if err := g.generateController(); err != nil {
		return err
	}

	// 生成路由
	if err := g.generateRoute(); err != nil {
		return err
	}

	// 生成前端代码
	if err := g.generateFrontend(); err != nil {
		return err
	}

	// 格式化生成的Go文件
	if err := g.FormatGoFiles(); err != nil {
		// 仅记录错误，不阻止生成过程
		fmt.Printf("警告：格式化Go文件时出错: %v\n", err)
	}

	// 记录生成历史
	_, err := g.History.Create(g.Config, g.generatedFiles)
	if err != nil {
		fmt.Printf("警告: 记录生成历史失败: %v\n", err)
	}

	return nil
}

// AddGeneratedFile 添加生成的文件路径
func (g *Generator) AddGeneratedFile(path, tmpl string) {
	// 将路径统一为相对项目根目录的路径
	relPath, err := filepath.Rel(g.RootPath, path)
	if err == nil {
		g.generatedFiles[relPath] = tmpl
	} else {
		g.generatedFiles[path] = tmpl
	}
}

// ListHistory 获取代码生成历史
func (g *Generator) ListHistory(page, pageSize int) ([]HistoryModel, int64, error) {
	return g.History.GetList(page, pageSize)
}

// RollBack 回滚代码生成
func (g *Generator) RollBack(id uint, deleteFiles, deleteAPI, deleteMenu, deleteTable bool) error {
	return g.History.RollBack(id, deleteFiles, deleteAPI, deleteMenu, deleteTable)
}

// ToLowerCamel 将字符串转换为小驼峰格式
func ToLowerCamel(s string) string {
	if s == "" {
		return s
	}
	result := []rune(s)
	result[0] = unicode.ToLower(result[0])
	return string(result)
}

// ToPlural 转换为复数形式 (简单实现)
func ToPlural(s string) string {
	if s == "" {
		return s
	}
	// 非常简单的转换，实际可能需要更复杂的规则
	if strings.HasSuffix(s, "y") {
		return s[:len(s)-1] + "ies"
	}
	return s + "s"
}

// FormatGoFiles 对所有生成的Go文件执行go fmt命令
func (g *Generator) FormatGoFiles() error {
	// 收集所有生成的.go文件
	var goFiles []string
	for filePath := range g.generatedFiles {
		if filepath.Ext(filePath) == ".go" {
			// 确保使用完整路径
			fullPath := filePath
			if !filepath.IsAbs(filePath) {
				fullPath = filepath.Join(g.RootPath, filePath)
			}
			goFiles = append(goFiles, fullPath)
		}
	}

	// 如果没有Go文件，则直接返回
	if len(goFiles) == 0 {
		return nil
	}

	// 对每个文件执行go fmt
	for _, file := range goFiles {
		fmt.Printf("执行go fmt: %s\n", file)
		cmd := exec.Command("go", "fmt", file)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("格式化文件 %s 失败: %w", file, err)
		}
	}

	return nil
}
