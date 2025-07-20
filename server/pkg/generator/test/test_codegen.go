package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/zhoudm1743/go-web/pkg/generator"
)

// 测试代码生成API
func main() {
	// 测试1: 创建基本模型
	testBasicModel()

	// 测试2: 创建带关系的模型
	testRelationModel()

	fmt.Println("测试完成！")
}

// testBasicModel 测试基本模型
func testBasicModel() {
	fmt.Println("1. 测试生成基本模型...")

	// 创建测试请求
	config := generator.Config{
		StructName:    "Product",
		TableName:     "products",
		PackageName:   "admin",
		Description:   "产品",
		ApiPrefix:     "products",
		AppName:       "admin",
		HasList:       true,
		HasCreate:     true,
		HasUpdate:     true,
		HasDelete:     true,
		HasDetail:     true,
		HasPagination: true,
		Fields: []*generator.Field{
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
				FieldDesc:    "产品名称",
				Required:     true,
				IsSearchable: true,
			},
			{
				FieldName:  "Price",
				FieldType:  "float64",
				ColumnName: "price",
				FieldDesc:  "价格",
				IsSortable: true,
			},
			{
				FieldName:    "Stock",
				FieldType:    "int",
				ColumnName:   "stock",
				FieldDesc:    "库存",
				IsFilterable: true,
			},
		},
	}

	// 发送请求
	sendGenerateRequest(config)
}

// testRelationModel 测试关系模型
func testRelationModel() {
	fmt.Println("2. 测试生成带关系的模型...")

	// 创建测试请求
	config := generator.Config{
		StructName:    "Article",
		TableName:     "articles",
		PackageName:   "admin",
		Description:   "文章",
		ApiPrefix:     "articles",
		AppName:       "admin",
		HasList:       true,
		HasCreate:     true,
		HasUpdate:     true,
		HasDelete:     true,
		HasDetail:     true,
		HasPagination: true,
		Fields: []*generator.Field{
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
			// 关系字段 - 从属于用户
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
			// 关系字段 - 拥有多个评论
			{
				FieldName:    "Comments",
				FieldDesc:    "评论列表",
				IsRelation:   true,
				RelationType: generator.HasMany,
				RelatedModel: "Comment",
				ForeignKey:   "ArticleID",
				References:   "ID",
			},
			// 关系字段 - 多对多标签
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
		},
	}

	// 发送请求
	sendGenerateRequest(config)
}

// sendGenerateRequest 发送代码生成请求
func sendGenerateRequest(config generator.Config) {
	// 打印请求内容
	fmt.Printf("请求生成 %s 模型...\n", config.StructName)

	// 将配置转为JSON
	jsonData, err := json.Marshal(config)
	if err != nil {
		fmt.Printf("JSON编码失败: %v\n", err)
		return
	}

	// 模拟HTTP请求 (实际应用中应调用真实的HTTP API)
	fmt.Println("模拟发送HTTP请求到 /admin/codegen/generate...")

	// 直接调用生成器
	gen := generator.New(&config)
	gen.SetRootPath("./test_output")

	// 初始化历史记录数据库
	if err := gen.InitHistoryDB(); err != nil {
		fmt.Printf("初始化历史记录数据库失败: %v\n", err)
		return
	}

	// 执行代码生成
	if err := gen.Run(); err != nil {
		fmt.Printf("生成代码失败: %v\n", err)
		return
	}

	fmt.Printf("%s 模型生成成功！\n", config.StructName)
	fmt.Println("生成的文件位于 ./test_output 目录")
	fmt.Println()

	// 等待1秒，便于观察
	time.Sleep(time.Second)
}

// TestHTTPRequest 测试HTTP请求 (实际应用中使用)
func TestHTTPRequest() {
	// 创建HTTP客户端
	client := &http.Client{}

	// 配置示例
	config := map[string]interface{}{
		"structName":    "Product",
		"tableName":     "products",
		"packageName":   "admin",
		"description":   "产品",
		"apiPrefix":     "products",
		"appName":       "admin",
		"hasList":       true,
		"hasCreate":     true,
		"hasUpdate":     true,
		"hasDelete":     true,
		"hasDetail":     true,
		"hasPagination": true,
		"fields": []map[string]interface{}{
			{
				"fieldName":    "Name",
				"fieldType":    "string",
				"columnName":   "name",
				"fieldDesc":    "产品名称",
				"required":     true,
				"isSearchable": true,
			},
		},
	}

	// 将配置转为JSON
	jsonData, _ := json.Marshal(config)

	// 创建请求
	req, err := http.NewRequest("POST", "http://localhost:8080/admin/codegen/generate", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("发送请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("响应状态码: %d\n", resp.StatusCode)
	fmt.Printf("响应内容: %s\n", string(body))
}
