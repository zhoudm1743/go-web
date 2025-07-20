package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-web/apps/admin/dto"
	"github.com/zhoudm1743/go-web/apps/admin/models"
	"github.com/zhoudm1743/go-web/core/facades"
	"github.com/zhoudm1743/go-web/core/response"
	"strconv"
)

// ProductController 产品控制器
type ProductController struct{}

// NewProductController 创建产品控制器
func NewProductController() *ProductController {
	return &ProductController{}
}

// GetProducts 获取产品列表
func (c *ProductController) GetProducts(ctx *gin.Context) {
	var params dto.ProductQueryParams
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
	var items []*models.Product
	
	query := db.Model(&models.Product{})
	
	// 应用查询条件
	
	if params.Name != "" {
		query = query.Where("name = ?", params.Name)
	}
	
	if params.Status != 0 {
		query = query.Where("status = ?", params.Status)
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
	result := &dto.ProductListResponse{
		Total: total,
		List:  make([]*dto.ProductResponse, len(items)),
	}

	for i, item := range items {
		resp := &dto.ProductResponse{
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



// GetProduct 获取产品详情
func (c *ProductController) GetProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	itemID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "ID无效")
		return
	}

	db := facades.DB()
	var item models.Product
	
	query := db
	

	if err := query.First(&item, itemID).Error; err != nil {
		response.FailWithMsg(ctx, response.Failed, "产品不存在")
		return
	}

	// 构造响应数据
	resp := &dto.ProductResponse{
		ID:        item.ID,
		CreatedAt: item.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: item.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	
	// 复制其他字段
	response.Copy(resp, item)

	response.OkWithData(ctx, resp)
}



// CreateProduct 创建产品
func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var req dto.ProductCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "请求参数有误")
		return
	}

	// 创建产品
	item := &models.Product{}
	response.Copy(item, req)

	db := facades.DB()
	if err := db.Create(item).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	response.OkWithMsg(ctx, "创建成功")
}



// UpdateProduct 更新产品
func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	var req dto.ProductUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "请求参数有误")
		return
	}

	db := facades.DB()
	var item models.Product
	if err := db.First(&item, req.ID).Error; err != nil {
		response.FailWithMsg(ctx, response.Failed, "产品不存在")
		return
	}

	// 只更新提供的字段
	updates := map[string]interface{}{}

	// 创建一个临时对象，用于复制非空字段
	tempProduct := &models.Product{}
	response.Copy(tempProduct, req)

	
	if req.Name != "" {
		updates["name"] = tempProduct.Name
	}
	
	if req.Status != 0 {
		updates["status"] = tempProduct.Status
	}
	
	if req.Price != 0 {
		updates["price"] = tempProduct.Price
	}
	

	if err := db.Model(&item).Updates(updates).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	response.OkWithMsg(ctx, "更新成功")
}



// DeleteProduct 删除产品
func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	itemID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "ID无效")
		return
	}

	db := facades.DB()
	if err := db.Delete(&models.Product{}, itemID).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	response.OkWithMsg(ctx, "删除成功")
}

