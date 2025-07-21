package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestAllFieldTypes 测试所有字段类型
func TestAllFieldTypes(t *testing.T) {
	// 创建临时测试目录
	testDir := filepath.Join(os.TempDir(), "go-web-generator-field-test")
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("无法创建测试目录: %v", err)
	}
	defer os.RemoveAll(testDir)

	// 创建必要的目录结构
	setupTestDirectories(t, testDir)

	// 创建基本框架文件
	setupBasicFiles(t, testDir)

	// 测试各种字段类型
	testBasicTypes(t, testDir)
	testRelationTypes(t, testDir)
}

// 创建必要的目录结构
func setupTestDirectories(t *testing.T, testDir string) {
	serverDir := filepath.Join(testDir, "server")
	appsDir := filepath.Join(serverDir, "apps")
	adminDir := filepath.Join(appsDir, "admin")
	modelsDir := filepath.Join(adminDir, "models")
	controllersDir := filepath.Join(adminDir, "controllers")
	dtoDir := filepath.Join(adminDir, "dto")
	routesDir := filepath.Join(adminDir, "routes")
	frontEndDir := filepath.Join(testDir, "front-end")
	serviceDir := filepath.Join(frontEndDir, "src/service/api")
	viewsDir := filepath.Join(frontEndDir, "src/views")

	// 创建所有必要目录
	for _, dir := range []string{modelsDir, controllersDir, dtoDir, routesDir, frontEndDir, serviceDir, viewsDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("无法创建目录 %s: %v", dir, err)
		}
	}
}

// 创建基本框架文件
func setupBasicFiles(t *testing.T, testDir string) {
	routeContent := `package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-web/apps/admin/controllers"
	"github.com/zhoudm1743/go-web/apps/admin/middlewares"
)

// InitRoutes 初始化路由
func RegisterRoutes(r *gin.RouterGroup) {
	// 初始化控制器
	authController := controllers.NewAuthController()
	
	publicRoutes := r
	{
		// 认证相关路由
		publicRoutes.POST("/login", authController.Login)
	}
	
	// 私有路由
	privateRoutes := r.Group("/admin")
	privateRoutes.Use(middlewares.AdminAuth())
	{
		// 认证相关路由
		privateRoutes.GET("/me", authController.GetUserInfo)
	}
}
`
	err := os.WriteFile(filepath.Join(testDir, "server/apps/admin/routes/routes.go"), []byte(routeContent), 0644)
	if err != nil {
		t.Fatalf("无法创建测试路由文件: %v", err)
	}
}

// 测试基本数据类型
func testBasicTypes(t *testing.T, testDir string) {
	config := &Config{
		StructName:  "BasicTypes",
		TableName:   "basic_types",
		PackageName: "admin",
		Description: "基本类型测试",
		ModuleName:  "github.com/zhoudm1743/go-web",
		RouterGroup: "privateRoutes",
		ApiPrefix:   "basicType",
		HasList:     true,
		HasCreate:   true,
		HasUpdate:   true,
		HasDelete:   true,
		HasDetail:   true,
		Fields: []*Field{
			{
				FieldName:    "ID",
				FieldType:    "uint",
				ColumnName:   "id",
				FieldDesc:    "主键ID",
				Required:     true,
				IsPrimaryKey: true,
				IsSearchable: true,
				IsFilterable: true,
				IsSortable:   true,
			},
			{
				FieldName:    "StringField",
				FieldType:    "string",
				ColumnName:   "string_field",
				FieldDesc:    "字符串字段",
				Required:     true,
				IsPrimaryKey: false,
				IsSearchable: true,
				IsFilterable: true,
				IsSortable:   false,
			},
			{
				FieldName:    "IntField",
				FieldType:    "int",
				ColumnName:   "int_field",
				FieldDesc:    "整数字段",
				Required:     false,
				IsPrimaryKey: false,
				IsSearchable: false,
				IsFilterable: true,
				IsSortable:   true,
			},
			{
				FieldName:    "Int64Field",
				FieldType:    "int64",
				ColumnName:   "int64_field",
				FieldDesc:    "64位整数字段",
				Required:     false,
				IsPrimaryKey: false,
				IsSearchable: false,
				IsFilterable: true,
				IsSortable:   true,
			},
			{
				FieldName:    "UintField",
				FieldType:    "uint",
				ColumnName:   "uint_field",
				FieldDesc:    "无符号整数字段",
				Required:     false,
				IsPrimaryKey: false,
				IsSearchable: false,
				IsFilterable: false,
				IsSortable:   true,
			},
			{
				FieldName:    "Uint64Field",
				FieldType:    "uint64",
				ColumnName:   "uint64_field",
				FieldDesc:    "64位无符号整数字段",
				Required:     false,
				IsPrimaryKey: false,
				IsSearchable: false,
				IsFilterable: false,
				IsSortable:   true,
			},
			{
				FieldName:    "Float64Field",
				FieldType:    "float64",
				ColumnName:   "float64_field",
				FieldDesc:    "浮点数字段",
				Required:     false,
				IsPrimaryKey: false,
				IsSearchable: false,
				IsFilterable: false,
				IsSortable:   true,
			},
			{
				FieldName:    "BoolField",
				FieldType:    "bool",
				ColumnName:   "bool_field",
				FieldDesc:    "布尔字段",
				Required:     false,
				IsPrimaryKey: false,
				IsSearchable: false,
				IsFilterable: true,
				IsSortable:   false,
			},
			{
				FieldName:    "TimeField",
				FieldType:    "time.Time",
				ColumnName:   "time_field",
				FieldDesc:    "时间字段",
				Required:     false,
				IsPrimaryKey: false,
				IsSearchable: false,
				IsFilterable: false,
				IsSortable:   true,
			},
		},
	}

	// 创建生成器
	gen := New(config)
	gen.SetRootPath(testDir)

	// 执行代码生成
	if err := gen.Run(); err != nil {
		t.Fatalf("代码生成失败: %v", err)
	}

	// 验证生成的模型文件
	modelFile := filepath.Join(testDir, "server/apps/admin/models/basicTypes.go")
	content, err := os.ReadFile(modelFile)
	if err != nil {
		t.Fatalf("读取模型文件失败: %v", err)
	}

	// 检查字段类型
	modelContent := string(content)
	expectedFields := map[string]string{
		"StringField":  "string",
		"IntField":     "int",
		"Int64Field":   "int64",
		"UintField":    "uint",
		"Uint64Field":  "uint64",
		"Float64Field": "float64",
		"BoolField":    "bool",
		"TimeField":    "time.Time",
	}

	for fieldName, fieldType := range expectedFields {
		if !strings.Contains(modelContent, fieldName+" "+fieldType) {
			t.Errorf("模型文件中缺少字段 %s %s", fieldName, fieldType)
		}
	}

	// 检查列名标签
	for fieldName, _ := range expectedFields {
		snakeCase := ToSnakeCase(fieldName)
		if !strings.Contains(modelContent, "gorm:\"column:"+snakeCase+"\"") {
			t.Errorf("模型文件中缺少正确的列名标签 gorm:\"column:%s\" 对于字段 %s", snakeCase, fieldName)
		}
	}
}

// 测试关系类型
func testRelationTypes(t *testing.T, testDir string) {
	config := &Config{
		StructName:  "RelationTypes",
		TableName:   "relation_types",
		PackageName: "admin",
		Description: "关系类型测试",
		ModuleName:  "github.com/zhoudm1743/go-web",
		RouterGroup: "privateRoutes",
		ApiPrefix:   "relationTypes",
		HasList:     true,
		HasCreate:   true,
		HasUpdate:   true,
		HasDelete:   true,
		HasDetail:   true,
		Fields: []*Field{
			{
				FieldName:    "ID",
				FieldType:    "uint",
				ColumnName:   "id",
				FieldDesc:    "主键ID",
				Required:     true,
				IsPrimaryKey: true,
				IsSearchable: true,
				IsFilterable: true,
				IsSortable:   true,
			},
			{
				FieldName:    "Name",
				FieldType:    "string",
				ColumnName:   "name",
				FieldDesc:    "名称",
				Required:     true,
				IsPrimaryKey: false,
				IsSearchable: true,
				IsFilterable: true,
				IsSortable:   false,
			},
			// BelongsTo 关系 - 默认命名
			{
				FieldName:    "Category",
				FieldType:    "string", // 会被覆盖
				ColumnName:   "",
				FieldDesc:    "分类",
				Required:     false,
				IsPrimaryKey: false,
				IsSearchable: false,
				IsFilterable: true,
				IsSortable:   false,
				IsRelation:   true,
				RelationType: "belongs_to",
				RelatedModel: "Category",
				ForeignKey:   "", // 使用默认命名 CategoryID
				References:   "", // 使用默认命名 ID
				Preload:      true,
				Joinable:     true,
			},
			// BelongsTo 关系 - 自定义命名
			{
				FieldName:    "Owner",
				FieldType:    "string", // 会被覆盖
				ColumnName:   "",
				FieldDesc:    "所有者",
				Required:     false,
				IsPrimaryKey: false,
				IsSearchable: false,
				IsFilterable: true,
				IsSortable:   false,
				IsRelation:   true,
				RelationType: "belongs_to",
				RelatedModel: "User",
				ForeignKey:   "OwnerID", // 自定义外键
				References:   "ID",
				Preload:      false,
				Joinable:     true,
			},
			// HasOne 关系
			{
				FieldName:    "Profile",
				FieldType:    "string", // 会被覆盖
				ColumnName:   "",
				FieldDesc:    "个人资料",
				Required:     false,
				IsPrimaryKey: false,
				IsSearchable: false,
				IsFilterable: false,
				IsSortable:   false,
				IsRelation:   true,
				RelationType: "has_one",
				RelatedModel: "Profile",
				ForeignKey:   "RelationTypesID", // 关联模型中的外键
				References:   "ID",
				Preload:      true,
				Joinable:     false,
			},
			// HasMany 关系
			{
				FieldName:    "Comments",
				FieldType:    "string", // 会被覆盖
				ColumnName:   "",
				FieldDesc:    "评论",
				Required:     false,
				IsPrimaryKey: false,
				IsSearchable: false,
				IsFilterable: false,
				IsSortable:   false,
				IsRelation:   true,
				RelationType: "has_many",
				RelatedModel: "Comment",
				ForeignKey:   "RelationTypesID", // 关联模型中的外键
				References:   "ID",
				Preload:      true,
				Joinable:     false,
			},
			// ManyToMany 关系
			{
				FieldName:    "Tags",
				FieldType:    "string", // 会被覆盖
				ColumnName:   "",
				FieldDesc:    "标签",
				Required:     false,
				IsPrimaryKey: false,
				IsSearchable: false,
				IsFilterable: false,
				IsSortable:   false,
				IsRelation:   true,
				RelationType: "many_to_many",
				RelatedModel: "Tag",
				JoinTable:    "relation_types_tags",
				ForeignKey:   "relation_types_id",
				References:   "tag_id",
				Preload:      true,
				Joinable:     false,
			},
		},
	}

	// 创建生成器
	gen := New(config)
	gen.SetRootPath(testDir)

	// 执行代码生成
	if err := gen.Run(); err != nil {
		t.Fatalf("代码生成失败: %v", err)
	}

	// 验证生成的模型文件
	modelFile := filepath.Join(testDir, "server/apps/admin/models/relationTypes.go")
	content, err := os.ReadFile(modelFile)
	if err != nil {
		t.Fatalf("读取模型文件失败: %v", err)
	}

	modelContent := string(content)

	// 验证关系字段
	relationFields := map[string][]string{
		// 字段名: [字段类型, 外键字段, 关系tag]
		"Category": {"*Category", "CategoryID uint", "foreignKey"},
		"Owner":    {"*User", "OwnerID uint", "foreignKey:OwnerID"},
		"Profile":  {"*Profile", "", "foreignKey:RelationTypesID"},
		"Comments": {"[]Comment", "", "foreignKey:RelationTypesID"},
		"Tags":     {"[]Tag", "", "many2many:relation_types_tags"},
	}

	for fieldName, checks := range relationFields {
		// 检查字段类型
		if len(checks) > 0 && checks[0] != "" {
			fieldType := checks[0]
			if !strings.Contains(modelContent, fieldName+" "+fieldType) {
				t.Errorf("模型文件中缺少正确的关系字段类型 %s %s", fieldName, fieldType)
			}
		}

		// 检查外键字段
		if len(checks) > 1 && checks[1] != "" {
			foreignKey := checks[1]
			if !strings.Contains(modelContent, foreignKey) {
				t.Errorf("模型文件中缺少外键字段 %s", foreignKey)
			}
		}

		// 检查关系标签
		if len(checks) > 2 && checks[2] != "" {
			relationTag := checks[2]
			if !strings.Contains(modelContent, relationTag) {
				t.Errorf("模型文件中缺少关系标签 %s 对于字段 %s", relationTag, fieldName)
			}
		}
	}

	// 检查预加载方法
	if !strings.Contains(modelContent, "LoadRelations") {
		t.Error("模型文件中缺少预加载方法 LoadRelations")
	}

	// 检查是否包含正确的Preload语句
	preloadFields := []string{"Category", "Profile", "Comments", "Tags"}
	for _, field := range preloadFields {
		preloadStmt := fmt.Sprintf("Preload(\"%s\")", field)
		if !strings.Contains(modelContent, preloadStmt) {
			t.Errorf("模型文件中缺少预加载语句 %s", preloadStmt)
		}
	}

	// 验证生成的控制器文件
	controllerFile := filepath.Join(testDir, "server/apps/admin/controllers/relationTypes_controller.go")
	content, err = os.ReadFile(controllerFile)
	if err != nil {
		t.Fatalf("读取控制器文件失败: %v", err)
	}

	controllerContent := string(content)

	// 检查JOIN语句
	if !strings.Contains(controllerContent, "Joins(") {
		t.Error("控制器文件中缺少JOIN语句")
	}

	// 验证生成的DTO文件
	dtoFile := filepath.Join(testDir, "server/apps/admin/dto/relationTypes.go")
	content, err = os.ReadFile(dtoFile)
	if err != nil {
		t.Fatalf("读取DTO文件失败: %v", err)
	}

	dtoContent := string(content)

	// 检查过滤字段
	filterFields := []string{"CategoryFilter", "OwnerFilter"}
	for _, field := range filterFields {
		if !strings.Contains(dtoContent, field+" string") {
			t.Errorf("DTO文件中缺少过滤字段 %s", field)
		}
	}
}
