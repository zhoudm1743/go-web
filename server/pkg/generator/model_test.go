package generator

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestModelGeneration 测试模型生成
func TestModelGeneration(t *testing.T) {
	// 创建临时测试目录
	testDir := filepath.Join(os.TempDir(), "go-web-generator-model-test")
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("无法创建测试目录: %v", err)
	}
	defer os.RemoveAll(testDir)

	// 创建必要的目录结构
	modelsDir := filepath.Join(testDir, "server/apps/admin/models")
	if err := os.MkdirAll(modelsDir, 0755); err != nil {
		t.Fatalf("无法创建目录 %s: %v", modelsDir, err)
	}

	// 测试BelongsTo关系
	t.Run("TestBelongsToRelation", func(t *testing.T) {
		config := &Config{
			StructName:  "Article",
			TableName:   "articles",
			PackageName: "admin",
			Description: "文章",
			ModuleName:  "github.com/zhoudm1743/go-web",
			RouterGroup: "privateRoutes",
			ApiPrefix:   "article",
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
				},
				{
					FieldName:    "Title",
					FieldType:    "string",
					ColumnName:   "title",
					FieldDesc:    "标题",
					Required:     true,
					IsPrimaryKey: false,
				},
				{
					FieldName:    "Content",
					FieldType:    "string",
					ColumnName:   "content",
					FieldDesc:    "内容",
					Required:     false,
					IsPrimaryKey: false,
				},
				{
					FieldName:    "Category",
					FieldType:    "string", // 会被覆盖
					ColumnName:   "",
					FieldDesc:    "分类",
					Required:     false,
					IsPrimaryKey: false,
					IsRelation:   true,
					RelationType: "belongs_to",
					RelatedModel: "Category",
					ForeignKey:   "", // 使用默认外键
					References:   "ID",
					Preload:      true,
				},
			},
		}

		// 创建生成器
		gen := New(config)
		gen.SetRootPath(testDir)

		// 仅生成模型
		err := gen.generateModel()
		if err != nil {
			t.Fatalf("生成模型失败: %v", err)
		}

		// 验证生成的模型文件
		modelFile := filepath.Join(testDir, "server/apps/admin/models/article.go")
		content, err := os.ReadFile(modelFile)
		if err != nil {
			t.Fatalf("读取模型文件失败: %v", err)
		}

		modelContent := string(content)
		t.Logf("生成的模型文件内容:\n%s", modelContent)

		// 检查是否包含正确的关系定义
		if !strings.Contains(modelContent, "Category *Category") {
			t.Error("模型文件中缺少正确的关联字段定义")
		}

		// 确保关联字段定义正确格式
		categoryPattern := regexp.MustCompile(`Category \*Category .*json:"category".*`)
		if !categoryPattern.MatchString(modelContent) {
			t.Error("模型文件中关联字段定义格式不正确")
		}

		// 检查是否包含正确的外键定义
		if !strings.Contains(modelContent, "CategoryID uint") {
			t.Error("模型文件中缺少正确的外键字段定义")
		}

		// 确保外键字段定义正确格式
		categoryIDPattern := regexp.MustCompile(`CategoryID uint .*column:category_id.*`)
		if !categoryIDPattern.MatchString(modelContent) {
			t.Error("模型文件中外键字段定义格式不正确")
		}

		// 确保外键的列名是蛇形命名
		if !strings.Contains(modelContent, `gorm:"column:category_id"`) {
			t.Error("模型文件中外键字段的列名不正确")
		}

		// 确保没有重复定义的字段
		categoryFieldCount := strings.Count(modelContent, "Category ")
		categoryIdFieldCount := strings.Count(modelContent, "CategoryID ")

		if categoryFieldCount != 1 {
			t.Errorf("Category字段被定义了%d次，期望1次", categoryFieldCount)
		}

		if categoryIdFieldCount != 1 {
			t.Errorf("CategoryID字段被定义了%d次，期望1次", categoryIdFieldCount)
		}
	})

	// 测试自定义外键名
	t.Run("TestCustomForeignKey", func(t *testing.T) {
		config := &Config{
			StructName:  "Post",
			TableName:   "posts",
			PackageName: "admin",
			Description: "帖子",
			ModuleName:  "github.com/zhoudm1743/go-web",
			RouterGroup: "privateRoutes",
			ApiPrefix:   "post",
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
				},
				{
					FieldName:    "Title",
					FieldType:    "string",
					ColumnName:   "title",
					FieldDesc:    "标题",
					Required:     true,
					IsPrimaryKey: false,
				},
				{
					FieldName:    "Author",
					FieldType:    "string", // 会被覆盖
					ColumnName:   "",
					FieldDesc:    "作者",
					Required:     false,
					IsPrimaryKey: false,
					IsRelation:   true,
					RelationType: "belongs_to",
					RelatedModel: "User",
					ForeignKey:   "AuthorID", // 自定义外键
					References:   "ID",
					Preload:      true,
				},
			},
		}

		// 创建生成器
		gen := New(config)
		gen.SetRootPath(testDir)

		// 仅生成模型
		err := gen.generateModel()
		if err != nil {
			t.Fatalf("生成模型失败: %v", err)
		}

		// 验证生成的模型文件
		modelFile := filepath.Join(testDir, "server/apps/admin/models/post.go")
		content, err := os.ReadFile(modelFile)
		if err != nil {
			t.Fatalf("读取模型文件失败: %v", err)
		}

		modelContent := string(content)
		t.Logf("生成的模型文件内容:\n%s", modelContent)

		// 检查是否包含正确的关系定义
		if !strings.Contains(modelContent, "Author *User") {
			t.Error("模型文件中缺少正确的关联字段定义")
		}

		// 检查是否包含自定义的外键定义
		if !strings.Contains(modelContent, "AuthorID uint") {
			t.Error("模型文件中缺少自定义外键字段定义")
		}

		// 确保外键的列名是蛇形命名
		if !strings.Contains(modelContent, `gorm:"column:author_id"`) {
			t.Error("模型文件中外键字段的列名不正确")
		}

		// 确保关系标签正确引用了外键
		if !strings.Contains(modelContent, `gorm:"foreignKey:AuthorID;references:ID"`) {
			t.Error("模型文件中关联字段的外键标签不正确")
		}
	})
}
