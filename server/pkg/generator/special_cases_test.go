package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestSpecialCases 测试特殊情况
func TestSpecialCases(t *testing.T) {
	// 创建临时测试目录
	testDir := filepath.Join(os.TempDir(), "go-web-generator-special-test")
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("无法创建测试目录: %v", err)
	}
	defer os.RemoveAll(testDir)

	// 创建必要的目录结构
	setupTestDirs(t, testDir)

	// 创建基本框架文件
	setupTestFiles(t, testDir)

	// 测试特殊命名情况
	testSpecialNaming(t, testDir)

	// 测试极端关系配置
	testComplexRelations(t, testDir)
}

func setupTestDirs(t *testing.T, testDir string) {
	dirs := []string{
		filepath.Join(testDir, "server/apps/admin/models"),
		filepath.Join(testDir, "server/apps/admin/controllers"),
		filepath.Join(testDir, "server/apps/admin/dto"),
		filepath.Join(testDir, "server/apps/admin/routes"),
		filepath.Join(testDir, "front-end/src/service/api"),
		filepath.Join(testDir, "front-end/src/views"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("无法创建目录 %s: %v", dir, err)
		}
	}
}

func setupTestFiles(t *testing.T, testDir string) {
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

// 测试特殊命名情况
func testSpecialNaming(t *testing.T, testDir string) {
	config := &Config{
		StructName:  "SpecialCase",
		TableName:   "special_cases",
		PackageName: "admin",
		Description: "特殊情况测试",
		ModuleName:  "github.com/zhoudm1743/go-web",
		RouterGroup: "privateRoutes",
		ApiPrefix:   "specialCase",
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
			// 特殊命名字段
			{
				FieldName:    "URLField",
				FieldType:    "string",
				ColumnName:   "url_field",
				FieldDesc:    "URL字段",
				Required:     true,
				IsPrimaryKey: false,
				IsSearchable: true,
				IsFilterable: true,
				IsSortable:   false,
			},
			// 字段名包含数字
			{
				FieldName:    "Field123",
				FieldType:    "string",
				ColumnName:   "field_123",
				FieldDesc:    "带数字的字段",
				Required:     false,
				IsPrimaryKey: false,
				IsSearchable: true,
				IsFilterable: false,
				IsSortable:   false,
			},
			// 空列名（应自动生成）
			{
				FieldName:    "EmptyColumnName",
				FieldType:    "string",
				ColumnName:   "",
				FieldDesc:    "空列名字段",
				Required:     false,
				IsPrimaryKey: false,
				IsSearchable: false,
				IsFilterable: false,
				IsSortable:   false,
			},
			// 特殊关系名称
			{
				FieldName:    "UserRole",
				FieldType:    "string", // 会被覆盖
				ColumnName:   "",
				FieldDesc:    "用户角色",
				Required:     false,
				IsPrimaryKey: false,
				IsSearchable: false,
				IsFilterable: true,
				IsSortable:   false,
				IsRelation:   true,
				RelationType: "belongs_to",
				RelatedModel: "Role",
				ForeignKey:   "RoleID", // 使用自定义外键
				References:   "ID",
				Preload:      true,
				Joinable:     true,
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
	modelFile := filepath.Join(testDir, "server/apps/admin/models/specialCase.go")
	content, err := os.ReadFile(modelFile)
	if err != nil {
		t.Fatalf("读取模型文件失败: %v", err)
	}

	modelContent := string(content)

	// 验证特殊命名字段
	if !strings.Contains(modelContent, "URLField string") {
		t.Error("模型文件中缺少正确的URL字段定义")
	}

	// 验证字段名包含数字
	if !strings.Contains(modelContent, "Field123 string") {
		t.Error("模型文件中缺少包含数字的字段定义")
	}

	// 验证空列名是否被正确处理
	if !strings.Contains(modelContent, `gorm:"column:empty_column_name"`) {
		t.Error("模型文件中空列名没有被正确处理")
	}

	// 验证特殊关系名称
	if !strings.Contains(modelContent, "UserRole *Role") {
		t.Error("模型文件中特殊关系字段名没有被正确处理")
	}

	// 验证DTO文件
	dtoFile := filepath.Join(testDir, "server/apps/admin/dto/specialcase.go")
	dtoContent, err := os.ReadFile(dtoFile)
	if err != nil {
		t.Fatalf("读取DTO文件失败: %v", err)
	}

	// 验证关联过滤字段
	if !strings.Contains(string(dtoContent), "UserRoleFilter string") {
		t.Error("DTO文件中缺少特殊命名的关联过滤字段")
	}
}

// 测试复杂关系配置
func testComplexRelations(t *testing.T, testDir string) {
	config := &Config{
		StructName:  "ComplexRelation",
		TableName:   "complex_relations",
		PackageName: "admin",
		Description: "复杂关系测试",
		ModuleName:  "github.com/zhoudm1743/go-web",
		RouterGroup: "privateRoutes",
		ApiPrefix:   "complexRelation",
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
			// 多个BelongsTo关系
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
				Preload:      true,
				Joinable:     true,
			},
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
				ForeignKey:   "OwnerID",
				References:   "ID",
				Preload:      true,
				Joinable:     true,
			},
			// 嵌套多层关系
			{
				FieldName:    "Team",
				FieldType:    "string", // 会被覆盖
				ColumnName:   "",
				FieldDesc:    "团队",
				Required:     false,
				IsPrimaryKey: false,
				IsSearchable: false,
				IsFilterable: true,
				IsSortable:   false,
				IsRelation:   true,
				RelationType: "belongs_to",
				RelatedModel: "Team",
				Preload:      true,
				Joinable:     true,
			},
			// 多对多关系
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
				JoinTable:    "complex_relation_tags",
				Preload:      true,
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
	modelFile := filepath.Join(testDir, "server/apps/admin/models/complexRelation.go")
	content, err := os.ReadFile(modelFile)
	if err != nil {
		t.Fatalf("读取模型文件失败: %v", err)
	}

	modelContent := string(content)

	// 验证多个BelongsTo关系
	expectedFields := map[string]string{
		"Category": "*Category",
		"Owner":    "*User",
		"Team":     "*Team",
		"Tags":     "[]Tag",
	}

	for fieldName, fieldType := range expectedFields {
		if !strings.Contains(modelContent, fieldName+" "+fieldType) {
			t.Errorf("模型文件中缺少关系字段 %s %s", fieldName, fieldType)
		}
	}

	// 验证外键字段
	expectedForeignKeys := []string{
		"CategoryID uint",
		"OwnerID uint",
		"TeamID uint",
	}

	for _, foreignKey := range expectedForeignKeys {
		if !strings.Contains(modelContent, foreignKey) {
			t.Errorf("模型文件中缺少外键字段 %s", foreignKey)
		}
	}

	// 验证关系标签
	relationTags := []string{
		`gorm:"many2many:complex_relation_tags"`,
		`gorm:"foreignKey:`,
	}

	for _, tag := range relationTags {
		if !strings.Contains(modelContent, tag) {
			t.Errorf("模型文件中缺少关系标签 %s", tag)
		}
	}

	// 验证预加载
	if !strings.Contains(modelContent, "LoadRelations") {
		t.Error("模型文件中缺少预加载方法")
	}

	// 验证DTO文件
	dtoFile := filepath.Join(testDir, "server/apps/admin/dto/complexrelation.go")
	dtoContent, err := os.ReadFile(dtoFile)
	if err != nil {
		t.Fatalf("读取DTO文件失败: %v", err)
	}

	// 验证关联过滤字段
	filterFields := []string{"CategoryFilter", "OwnerFilter", "TeamFilter"}
	for _, field := range filterFields {
		if !strings.Contains(string(dtoContent), field+" string") {
			t.Errorf("DTO文件中缺少关联过滤字段 %s", field)
		}
	}

	// 验证控制器文件
	controllerFile := filepath.Join(testDir, "server/apps/admin/controllers/complexrelation_controller.go")
	controllerContent, err := os.ReadFile(controllerFile)
	if err != nil {
		t.Fatalf("读取控制器文件失败: %v", err)
	}

	// 验证JOIN和Preload
	ctrlStr := string(controllerContent)
	if !strings.Contains(ctrlStr, "Joins(") {
		t.Error("控制器文件中缺少JOIN语句")
	}

	preloadFields := []string{"Category", "Owner", "Team", "Tags"}
	for _, field := range preloadFields {
		if !strings.Contains(ctrlStr, "Preload(\""+field+"\")") {
			t.Errorf("控制器文件中缺少对 %s 的预加载", field)
		}
	}
}
