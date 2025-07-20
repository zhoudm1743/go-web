package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// Field 字段定义
type Field struct {
	FieldName    string // 结构体字段名称
	FieldType    string // Go类型
	ColumnName   string // 数据库字段名
	FieldDesc    string // 字段描述
	Required     bool   // 是否必填
	IsPrimaryKey bool   // 是否主键
	IsSearchable bool   // 是否可搜索
	IsFilterable bool   // 是否可过滤
	IsSortable   bool   // 是否可排序

	// 关系字段
	IsRelation   bool   // 是否关系字段
	RelationType string // 关系类型
	RelatedModel string // 关联模型
	ForeignKey   string // 外键 (本模型字段)
	References   string // 引用字段 (关联模型字段)
	Preload      bool   // 是否预加载
	JoinTable    string // 多对多关联表名
}

func main() {
	// 创建测试目录
	os.MkdirAll("./test_output", 0755)

	// 测试生成带关系的模型
	fmt.Println("测试生成带关系的模型...")

	// 生成模型文件
	modelFile := filepath.Join("./test_output", "article.go")
	if err := os.WriteFile(modelFile, []byte(generateModelCode()), 0644); err != nil {
		fmt.Printf("写入模型文件失败: %v\n", err)
		return
	}

	// 生成DTO文件
	dtoFile := filepath.Join("./test_output", "article_dto.go")
	if err := os.WriteFile(dtoFile, []byte(generateDTOCode()), 0644); err != nil {
		fmt.Printf("写入DTO文件失败: %v\n", err)
		return
	}

	// 生成控制器文件
	controllerFile := filepath.Join("./test_output", "article_controller.go")
	if err := os.WriteFile(controllerFile, []byte(generateControllerCode()), 0644); err != nil {
		fmt.Printf("写入控制器文件失败: %v\n", err)
		return
	}

	fmt.Println("测试完成，请查看test_output目录中的生成文件。")
}

// generateModelCode 生成模型代码
func generateModelCode() string {
	return `package models

import (
	"time"
	"github.com/zhoudm1743/go-web/core/facades"
	"gorm.io/gorm"
)

// Article 文章
type Article struct {
	ID        uint           ` + "`" + `gorm:"primarykey" json:"id"` + "`" + `                // 主键ID
	CreatedAt time.Time      ` + "`" + `json:"createdAt"` + "`" + `                           // 创建时间
	UpdatedAt time.Time      ` + "`" + `json:"updatedAt"` + "`" + `                           // 更新时间
	DeletedAt gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `                      // 删除时间
	Title     string         ` + "`" + `gorm:"column:title" json:"title"` + "`" + `           // 标题
	Content   string         ` + "`" + `gorm:"column:content" json:"content"` + "`" + `       // 内容
	UserID    uint           ` + "`" + `gorm:"column:user_id" json:"userId"` + "`" + `        // User外键
	User      User           ` + "`" + `gorm:"foreignKey:UserID;references:ID" json:"user"` + "`" + `  // 用户
	Comments  []Comment      ` + "`" + `gorm:"foreignKey:ArticleID" json:"comments"` + "`" + `        // 评论列表
	Tags      []Tag          ` + "`" + `gorm:"many2many:article_tags;foreignKey:article_id;references:tag_id" json:"tags"` + "`" + ` // 标签列表
}

// TableName 设置表名
func (Article) TableName() string {
	return "articles"
}

// 关系预加载
func (m *Article) LoadRelations(db *gorm.DB) *gorm.DB {
	query := db
	query = query.Preload("User")
	query = query.Preload("Tags")
	return query
}
`
}

// generateDTOCode 生成DTO代码
func generateDTOCode() string {
	return `package dto

import "time"

// ArticleCreateRequest 创建文章请求
type ArticleCreateRequest struct {
	Title   string ` + "`" + `json:"title" binding:"required"` + "`" + ` // 标题
	Content string ` + "`" + `json:"content"` + "`" + ` // 内容
	UserID  uint   ` + "`" + `json:"userId" binding:"required"` + "`" + ` // 用户 ID
}

// ArticleUpdateRequest 更新文章请求
type ArticleUpdateRequest struct {
	ID      uint   ` + "`" + `json:"id" binding:"required"` + "`" + ` // ID
	Title   string ` + "`" + `json:"title"` + "`" + ` // 标题
	Content string ` + "`" + `json:"content"` + "`" + ` // 内容
	UserID  uint   ` + "`" + `json:"userId"` + "`" + ` // 用户 ID
}

// ArticleQueryParams 文章查询参数
type ArticleQueryParams struct {
	Page          int  ` + "`" + `form:"page"` + "`" + `      // 页码
	PageSize      int  ` + "`" + `form:"pageSize"` + "`" + `  // 每页条数
	WithRelations bool ` + "`" + `form:"withRelations"` + "`" + ` // 是否加载关联
	Title         string ` + "`" + `form:"title"` + "`" + ` // 标题
	UserID        uint   ` + "`" + `form:"userId"` + "`" + ` // 用户 ID
}

// ArticleResponse 文章响应
type ArticleResponse struct {
	ID        uint       ` + "`" + `json:"id"` + "`" + `         // ID
	CreatedAt string     ` + "`" + `json:"createdAt"` + "`" + `  // 创建时间
	UpdatedAt string     ` + "`" + `json:"updatedAt"` + "`" + `  // 更新时间
	Title     string     ` + "`" + `json:"title"` + "`" + `      // 标题
	Content   string     ` + "`" + `json:"content"` + "`" + `    // 内容
	User      *UserResponse ` + "`" + `json:"user"` + "`" + `    // 用户
	Comments  []*CommentResponse ` + "`" + `json:"comments"` + "`" + ` // 评论列表
	Tags      []*TagResponse ` + "`" + `json:"tags"` + "`" + `   // 标签列表
}

// ArticleListResponse 文章列表响应
type ArticleListResponse struct {
	Total int64             ` + "`" + `json:"total"` + "`" + `  // 总数
	List  []*ArticleResponse ` + "`" + `json:"list"` + "`" + `  // 列表
}

// UserResponse 用户响应
type UserResponse struct {
	ID   uint   ` + "`" + `json:"id"` + "`" + `   // ID
	Name string ` + "`" + `json:"name"` + "`" + ` // 名称
}

// CommentResponse 评论响应
type CommentResponse struct {
	ID      uint   ` + "`" + `json:"id"` + "`" + `      // ID
	Content string ` + "`" + `json:"content"` + "`" + ` // 内容
}

// TagResponse 标签响应
type TagResponse struct {
	ID   uint   ` + "`" + `json:"id"` + "`" + `   // ID
	Name string ` + "`" + `json:"name"` + "`" + ` // 名称
}
`
}

// generateControllerCode 生成控制器代码
func generateControllerCode() string {
	return `package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-web/apps/admin/dto"
	"github.com/zhoudm1743/go-web/apps/admin/models"
	"github.com/zhoudm1743/go-web/core/facades"
	"github.com/zhoudm1743/go-web/core/response"
	"strconv"
)

// ArticleController 文章控制器
type ArticleController struct{}

// NewArticleController 创建文章控制器
func NewArticleController() *ArticleController {
	return &ArticleController{}
}

// GetArticles 获取文章列表
func (c *ArticleController) GetArticles(ctx *gin.Context) {
	var params dto.ArticleQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "请求参数有误")
		return
	}

	// 默认分页参数
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	db := facades.DB()
	var total int64
	var items []*models.Article
	
	query := db.Model(&models.Article{})
	
	// 应用查询条件
	if params.Title != "" {
		query = query.Where("title LIKE ?", "%"+params.Title+"%")
	}
	if params.UserID != 0 {
		query = query.Where("user_id = ?", params.UserID)
	}

	// 应用预加载
	if params.WithRelations {
		query = (&models.Article{}).LoadRelations(query)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	// 分页查询
	if err := query.Offset((params.Page - 1) * params.PageSize).
		Limit(params.PageSize).
		Find(&items).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	// 构造响应数据
	result := &dto.ArticleListResponse{
		Total: total,
		List:  make([]*dto.ArticleResponse, len(items)),
	}

	for i, item := range items {
		resp := &dto.ArticleResponse{
			ID:        item.ID,
			CreatedAt: item.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: item.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		
		// 复制其他字段
		response.Copy(resp, item)
		result.List[i] = resp
	}

	response.OkWithData(ctx, result)
}

// GetArticle 获取文章详情
func (c *ArticleController) GetArticle(ctx *gin.Context) {
	id := ctx.Param("id")
	itemID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "ID无效")
		return
	}

	db := facades.DB()
	var item models.Article
	
	query := db
	// 预加载关联数据
	withRelations := ctx.Query("withRelations")
	if withRelations == "true" {
		query = (&models.Article{}).LoadRelations(query)
	}

	if err := query.First(&item, itemID).Error; err != nil {
		response.FailWithMsg(ctx, response.Failed, "文章不存在")
		return
	}

	// 构造响应数据
	resp := &dto.ArticleResponse{
		ID:        item.ID,
		CreatedAt: item.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: item.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	
	// 复制其他字段
	response.Copy(resp, item)

	response.OkWithData(ctx, resp)
}

// CreateArticle 创建文章
func (c *ArticleController) CreateArticle(ctx *gin.Context) {
	var req dto.ArticleCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "请求参数有误")
		return
	}

	// 创建文章
	item := &models.Article{}
	response.Copy(item, req)

	db := facades.DB()
	if err := db.Create(item).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	response.OkWithMsg(ctx, "创建成功")
}

// UpdateArticle 更新文章
func (c *ArticleController) UpdateArticle(ctx *gin.Context) {
	var req dto.ArticleUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "请求参数有误")
		return
	}

	db := facades.DB()
	var item models.Article
	if err := db.First(&item, req.ID).Error; err != nil {
		response.FailWithMsg(ctx, response.Failed, "文章不存在")
		return
	}

	// 只更新提供的字段
	updates := map[string]interface{}{}

	// 创建一个临时对象，用于复制非空字段
	tempArticle := &models.Article{}
	response.Copy(tempArticle, req)

	if req.Title != "" {
		updates["title"] = tempArticle.Title
	}
	if req.Content != "" {
		updates["content"] = tempArticle.Content
	}
	if req.UserID != 0 {
		updates["user_id"] = tempArticle.UserID
	}

	if err := db.Model(&item).Updates(updates).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	response.OkWithMsg(ctx, "更新成功")
}

// DeleteArticle 删除文章
func (c *ArticleController) DeleteArticle(ctx *gin.Context) {
	id := ctx.Param("id")
	itemID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "ID无效")
		return
	}

	db := facades.DB()
	if err := db.Delete(&models.Article{}, itemID).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	response.OkWithMsg(ctx, "删除成功")
}`
}
