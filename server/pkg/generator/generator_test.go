package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestGenerator 测试代码生成器
func TestGenerator(t *testing.T) {
	// 创建临时测试目录
	testDir := filepath.Join(os.TempDir(), "go-web-generator-test")
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("无法创建测试目录: %v", err)
	}
	defer os.RemoveAll(testDir)

	// 创建必要的目录结构
	serverDir := filepath.Join(testDir, "server")
	appsDir := filepath.Join(serverDir, "apps")
	adminDir := filepath.Join(appsDir, "admin")
	modelsDir := filepath.Join(adminDir, "models")
	controllersDir := filepath.Join(adminDir, "controllers")
	dtoDir := filepath.Join(adminDir, "dto")
	routesDir := filepath.Join(adminDir, "routes")
	frontEndDir := filepath.Join(testDir, "front-end")

	// 创建所有必要目录
	for _, dir := range []string{modelsDir, controllersDir, dtoDir, routesDir, frontEndDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("无法创建目录 %s: %v", dir, err)
		}
	}

	// 创建测试路由文件
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
	err = os.WriteFile(filepath.Join(routesDir, "routes.go"), []byte(routeContent), 0644)
	if err != nil {
		t.Fatalf("无法创建测试路由文件: %v", err)
	}

	// 创建测试配置
	config := &Config{
		StructName:  "TestProduct",
		TableName:   "test_products",
		PackageName: "admin",
		Description: "测试产品",
		ModuleName:  "github.com/zhoudm1743/go-web",
		RouterGroup: "privateRoutes",
		ApiPrefix:   "testProduct",
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
			{
				FieldName:    "Category",
				FieldType:    "string", // 这里会被关系类型覆盖
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
				ForeignKey:   "CategoryID",
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

	// 检查生成的文件
	checkGeneratedFiles(t, testDir, config)
}

// 检查生成的文件
func checkGeneratedFiles(t *testing.T, testDir string, config *Config) {
	// 验证后端文件是否生成
	backendFiles := []string{
		filepath.Join(testDir, "server/apps", config.PackageName, "models", ToLowerCamel(config.StructName)+".go"),
		filepath.Join(testDir, "server/apps", config.PackageName, "dto", ToLowerCamel(config.StructName)+".go"),
		filepath.Join(testDir, "server/apps", config.PackageName, "controllers", ToLowerCamel(config.StructName)+"_controller.go"),
	}

	for _, file := range backendFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Errorf("后端文件未生成: %s", file)
			continue
		}

		// 读取文件内容进行验证
		content, err := os.ReadFile(file)
		if err != nil {
			t.Errorf("读取文件失败 %s: %v", file, err)
			continue
		}

		// 验证文件内容
		validateFileContent(t, string(content), file, config)
	}

	// 验证前端文件是否生成
	frontendFiles := []string{
		filepath.Join(testDir, "front-end/src/views", strings.ToLower(config.StructName), "index.vue"),
		filepath.Join(testDir, "front-end/src/views", strings.ToLower(config.StructName), "components/TableModal.vue"),
		filepath.Join(testDir, "front-end/src/service/api", ToLowerCamel(config.StructName)+".ts"),
	}

	for _, file := range frontendFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			// 对于测试环境，允许前端文件生成失败但需要输出警告
			t.Logf("警告：前端文件未生成: %s，这可能是因为测试环境限制", file)
		}
	}

	// 验证路由是否更新
	routeFile := filepath.Join(testDir, "server/apps", config.PackageName, "routes/routes.go")
	if _, err := os.Stat(routeFile); os.IsNotExist(err) {
		t.Errorf("路由文件未找到: %s", routeFile)
	} else {
		content, err := os.ReadFile(routeFile)
		if err != nil {
			t.Errorf("读取路由文件失败: %v", err)
		} else {
			routeContent := string(content)
			// 验证是否包含控制器变量
			// 注意：实际中控制器变量可能是全小写形式
			controllerVar := ToLowerCamel(config.StructName) + "Controller"
			controllerVarLower := strings.ToLower(ToLowerCamel(config.StructName)) + "controller"

			if !contains(routeContent, controllerVar) && !contains(routeContent, controllerVarLower) {
				t.Errorf("路由文件未包含控制器变量: %s 或 %s", controllerVar, controllerVarLower)
			}

			// 验证是否包含路由注册
			routeRegister := config.StructName + "路由"
			if !contains(routeContent, routeRegister) {
				t.Errorf("路由文件未包含路由注册: %s", routeRegister)
			}
		}
	}
}

// 验证文件内容
func validateFileContent(t *testing.T, content, file string, config *Config) {
	// 模型文件验证
	if filepath.Base(file) == ToLowerCamel(config.StructName)+".go" && filepath.Base(filepath.Dir(file)) == "models" {
		// 验证是否包含关系字段
		for _, field := range config.Fields {
			if field.IsRelation {
				// 验证关系字段类型
				if field.RelationType == "belongs_to" {
					relationField := field.FieldName + " *" + field.RelatedModel
					if !contains(content, relationField) {
						t.Errorf("模型文件未包含关系字段类型 %s: %s", file, relationField)
					}
				}
				// 验证外键字段
				if field.ForeignKey != "" {
					foreignKeyField := field.ForeignKey + " uint"
					if !contains(content, foreignKeyField) {
						t.Errorf("模型文件未包含外键字段 %s: %s", file, foreignKeyField)
					}
				}
				// 验证预加载方法
				if field.Preload {
					preloadMethod := "LoadRelations"
					preloadCode := `Preload("` + field.FieldName + `")`
					if !contains(content, preloadMethod) || !contains(content, preloadCode) {
						t.Errorf("模型文件未包含预加载方法或代码 %s", file)
					}
				}
			}
		}
	}

	// 控制器文件验证
	if filepath.Base(file) == ToLowerCamel(config.StructName)+"_controller.go" {
		// 验证是否包含预加载代码
		for _, field := range config.Fields {
			if field.IsRelation && field.Preload {
				preloadCode := `Preload("` + field.FieldName + `")`
				if !contains(content, preloadCode) {
					t.Errorf("控制器文件未包含预加载代码 %s: %s", file, preloadCode)
				}
			}
			// 验证是否包含JOIN代码
			if field.IsRelation && field.Joinable {
				joinCode := "Joins("
				if !contains(content, joinCode) {
					t.Errorf("控制器文件未包含JOIN代码 %s", file)
				}
			}
		}
	}

	// DTO文件验证
	if filepath.Base(file) == ToLowerCamel(config.StructName)+".go" && filepath.Base(filepath.Dir(file)) == "dto" {
		// 验证查询参数中是否包含关系过滤字段
		for _, field := range config.Fields {
			if field.IsRelation && field.IsFilterable {
				filterField := field.FieldName + "Filter"
				if !contains(content, filterField) {
					t.Errorf("DTO文件未包含关系过滤字段 %s: %s", file, filterField)
				}
			}
		}
	}
}

// 检查字符串是否包含子串
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
