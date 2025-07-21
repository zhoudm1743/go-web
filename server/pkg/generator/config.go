package generator

import (
	"strings"
)

// Config 代码生成器配置
type Config struct {
	// 基本信息
	StructName  string // 结构体名称
	TableName   string // 表名
	PackageName string // 包名
	Description string // 描述
	ModuleName  string // 模块名(import路径前缀)

	// 路由和API
	RouterGroup string // 路由分组
	ApiPrefix   string // API前缀

	// 选项
	HasList       bool // 是否有列表
	HasCreate     bool // 是否有创建
	HasUpdate     bool // 是否有更新
	HasDelete     bool // 是否有删除
	HasDetail     bool // 是否有详情
	HasPagination bool // 是否分页

	// 字段配置
	Fields []*Field
}

// Field 字段配置
type Field struct {
	FieldName    string // 结构体字段名称
	FieldType    string // Go类型
	ColumnName   string // 数据库字段名
	ColumnType   string // 数据库字段类型
	FieldDesc    string // 字段描述
	Required     bool   // 是否必填
	IsPrimaryKey bool   // 是否主键
	IsSearchable bool   // 是否可搜索
	IsFilterable bool   // 是否可过滤
	IsSortable   bool   // 是否可排序

	// 关系字段
	IsRelation   bool         // 是否关系字段
	RelationType RelationType // 关系类型
	RelatedModel string       // 关联模型
	ForeignKey   string       // 外键 (本模型字段)
	References   string       // 引用字段 (关联模型字段)
	Preload      bool         // 是否预加载
	JoinTable    string       // 多对多关联表名

	// JOIN查询相关
	Joinable        bool   // 是否支持JOIN查询
	JoinCondition   string // JOIN条件字段
	FilterCondition string // 过滤条件字段
}

// RelationType 关系类型枚举
type RelationType string

const (
	BelongsTo  RelationType = "belongs_to"   // 从属于
	HasOne     RelationType = "has_one"      // 拥有一个
	HasMany    RelationType = "has_many"     // 拥有多个
	ManyToMany RelationType = "many_to_many" // 多对多
)

// DatabaseInfo 数据库信息
type DatabaseInfo struct {
	Database string `json:"database"`
}

// TableInfo 表信息
type TableInfo struct {
	TableName    string `json:"tableName"`
	TableComment string `json:"tableComment"`
}

// ColumnInfo 列信息
type ColumnInfo struct {
	ColumnName    string `json:"columnName"`
	DataType      string `json:"dataType"`
	ColumnComment string `json:"columnComment"`
	IsNullable    string `json:"isNullable"`
	ColumnKey     string `json:"columnKey"`
}

// ConvertDataTypeToGo 将数据库类型转换为Go类型
func ConvertDataTypeToGo(dataType string) string {
	dataType = strings.ToLower(dataType)

	switch {
	// Boolean类型
	case dataType == "tinyint(1)" || dataType == "boolean" || dataType == "bool":
		return "bool"

	// 整数类型
	case strings.Contains(dataType, "int") && !strings.Contains(dataType, "point"):
		if strings.Contains(dataType, "big") || strings.Contains(dataType, "int8") {
			return "int64"
		}
		return "int"

	// 浮点类型
	case dataType == "float" || dataType == "double" || dataType == "decimal" ||
		dataType == "numeric" || dataType == "real" || strings.Contains(dataType, "float"):
		return "float64"

	// 字符串类型
	case dataType == "char" || dataType == "varchar" || dataType == "text" ||
		dataType == "mediumtext" || dataType == "longtext" ||
		dataType == "character varying" || dataType == "character" ||
		strings.Contains(dataType, "char") || strings.Contains(dataType, "text"):
		return "string"

	// 时间类型
	case dataType == "date" || dataType == "datetime" || dataType == "timestamp" ||
		strings.Contains(dataType, "time") || strings.Contains(dataType, "date"):
		return "time.Time"

	// JSON类型
	case dataType == "json" || dataType == "jsonb":
		return "json.RawMessage"

	// 二进制类型
	case dataType == "blob" || dataType == "bytea" || dataType == "binary" ||
		dataType == "varbinary" || strings.Contains(dataType, "blob") ||
		strings.Contains(dataType, "binary"):
		return "[]byte"

	// 默认使用string
	default:
		return "string"
	}
}
