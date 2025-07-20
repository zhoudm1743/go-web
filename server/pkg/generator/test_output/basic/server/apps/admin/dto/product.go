package dto



// ProductCreateRequest 创建产品请求
type ProductCreateRequest struct {

	Name string `json:"name" binding:"required"` // 名称

	Status int `json:"status" ` // 状态

	Price float64 `json:"price" ` // 价格

}

// ProductUpdateRequest 更新产品请求
type ProductUpdateRequest struct {
	ID uint `json:"id" binding:"required"` // ID

	Name string `json:"name" ` // 名称

	Status int `json:"status" ` // 状态

	Price float64 `json:"price" ` // 价格

}

// ProductQueryParams 产品查询参数
type ProductQueryParams struct {
	Page     int `form:"page"`      // 页码
	PageSize int `form:"pageSize"`  // 每页条数


	Name string `form:"name"` // 名称

	Status int `form:"status"` // 状态

}

// ProductResponse 产品响应
type ProductResponse struct {
	ID        uint   `json:"id"`         // ID
	CreatedAt string `json:"createdAt"`  // 创建时间
	UpdatedAt string `json:"updatedAt"`  // 更新时间

	Name string `json:"name"` // 名称

	Status int `json:"status"` // 状态

	Price float64 `json:"price"` // 价格

}

// ProductListResponse 产品列表响应
type ProductListResponse struct {
	Total int64                  `json:"total"`  // 总数
	List  []*ProductResponse `json:"list"`    // 列表
}
