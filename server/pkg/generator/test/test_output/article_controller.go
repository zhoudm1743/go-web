package controllers

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
}