package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/zhoudm1743/go-web/pkg/generator"
)

func main() {
	// 确保目录存在
	os.MkdirAll("./test_output", 0755)

	// 测试1：基本模型
	testBasicModel()

	// 测试2：关系模型
	testRelationModel()

	fmt.Println("测试完成，请查看test_output目录")
}

// testBasicModel 测试基本模型生成
func testBasicModel() {
	fmt.Println("开始测试基本模型生成...")

	// 创建基本字段
	fields := []*generator.Field{
		{
			FieldName:    "ID",
			FieldType:    "uint",
			ColumnName:   "id",
			FieldDesc:    "主键ID",
			IsPrimaryKey: true,
		},
		{
			FieldName:    "Name",
			FieldType:    "string",
			ColumnName:   "name",
			FieldDesc:    "名称",
			Required:     true,
			IsSearchable: true,
		},
		{
			FieldName:    "Status",
			FieldType:    "int",
			ColumnName:   "status",
			FieldDesc:    "状态",
			IsFilterable: true,
		},
		{
			FieldName:  "Price",
			FieldType:  "float64",
			ColumnName: "price",
			FieldDesc:  "价格",
			IsSortable: true,
		},
	}

	// 创建配置
	config := &generator.Config{
		StructName:    "Product",
		TableName:     "products",
		PackageName:   "admin",
		Description:   "产品",
		ModuleName:    "github.com/zhoudm1743/go-web",
		RouterGroup:   "privateRoutes",
		ApiPrefix:     "products",
		HasList:       true,
		HasCreate:     true,
		HasUpdate:     true,
		HasDelete:     true,
		HasDetail:     true,
		HasPagination: true,
		Fields:        fields,
	}

	// 创建生成器
	gen := generator.New(config)
	gen.SetRootPath("./test_output/basic")

	// 初始化目录结构
	setupTestDirs("./test_output/basic")

	// 执行生成
	if err := gen.Run(); err != nil {
		fmt.Printf("生成基本模型失败: %v\n", err)
		return
	}

	fmt.Println("基本模型生成成功!")
}

// testRelationModel 测试关系模型生成
func testRelationModel() {
	fmt.Println("开始测试关系模型生成...")

	// 创建字段，包含关系字段
	fields := []*generator.Field{
		{
			FieldName:    "ID",
			FieldType:    "uint",
			ColumnName:   "id",
			FieldDesc:    "主键ID",
			IsPrimaryKey: true,
		},
		{
			FieldName:    "Title",
			FieldType:    "string",
			ColumnName:   "title",
			FieldDesc:    "标题",
			Required:     true,
			IsSearchable: true,
		},
		{
			FieldName:  "Content",
			FieldType:  "string",
			ColumnName: "content",
			FieldDesc:  "内容",
		},
		// 关系字段 - 多对一，文章属于用户
		{
			FieldName:    "User",
			FieldDesc:    "用户",
			IsRelation:   true,
			RelationType: generator.BelongsTo,
			RelatedModel: "User",
			ForeignKey:   "UserID",
			References:   "ID",
			Preload:      true,
		},
		// 关系字段 - 一对多，文章有多个评论
		{
			FieldName:    "Comments",
			FieldDesc:    "评论列表",
			IsRelation:   true,
			RelationType: generator.HasMany,
			RelatedModel: "Comment",
			ForeignKey:   "ArticleID",
			References:   "ID",
		},
		// 关系字段 - 多对多，文章有多个标签
		{
			FieldName:    "Tags",
			FieldDesc:    "标签列表",
			IsRelation:   true,
			RelationType: generator.ManyToMany,
			RelatedModel: "Tag",
			JoinTable:    "article_tags",
			ForeignKey:   "article_id",
			References:   "tag_id",
			Preload:      true,
		},
	}

	// 创建配置
	config := &generator.Config{
		StructName:    "Article",
		TableName:     "articles",
		PackageName:   "admin",
		Description:   "文章",
		ModuleName:    "github.com/zhoudm1743/go-web",
		RouterGroup:   "privateRoutes",
		ApiPrefix:     "articles",
		HasList:       true,
		HasCreate:     true,
		HasUpdate:     true,
		HasDelete:     true,
		HasDetail:     true,
		HasPagination: true,
		Fields:        fields,
	}

	// 创建生成器
	gen := generator.New(config)
	gen.SetRootPath("./test_output/relation")

	// 初始化目录结构
	setupTestDirs("./test_output/relation")

	// 执行生成
	if err := gen.Run(); err != nil {
		fmt.Printf("生成关系模型失败: %v\n", err)
		return
	}

	fmt.Println("关系模型生成成功!")
}

// setupTestDirs 设置测试目录结构
func setupTestDirs(rootPath string) {
	// 创建目录结构
	dirs := []string{
		filepath.Join(rootPath, "server/apps/admin/controllers"),
		filepath.Join(rootPath, "server/apps/admin/models"),
		filepath.Join(rootPath, "server/apps/admin/dto"),
		filepath.Join(rootPath, "server/apps/admin/routes"),
		filepath.Join(rootPath, "server/temp/snippets"),
		filepath.Join(rootPath, "front-end/src/views/admin"),
		filepath.Join(rootPath, "front-end/src/service/api"),
	}

	// 创建目录
	for _, dir := range dirs {
		os.MkdirAll(dir, 0755)
	}

	// 创建路由文件
	routesFile := filepath.Join(rootPath, "server/apps/admin/routes/routes.go")
	routesContent := `package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-web/apps/admin/controllers"
)

// InitRoutes 初始化路由
func InitRoutes(r *gin.Engine) {
	// 初始化控制器

	// 私有路由
	privateRoutes := r.Group("/admin")
	{
		// 路由示例
		// privateRoutes.GET("/example", exampleController.Example)
	}
}`

	os.WriteFile(routesFile, []byte(routesContent), 0644)
}
