package dto

import "time"

// ArticleCreateRequest 创建文章请求
type ArticleCreateRequest struct {
	Title   string `json:"title" binding:"required"` // 标题
	Content string `json:"content"` // 内容
	UserID  uint   `json:"userId" binding:"required"` // 用户 ID
}

// ArticleUpdateRequest 更新文章请求
type ArticleUpdateRequest struct {
	ID      uint   `json:"id" binding:"required"` // ID
	Title   string `json:"title"` // 标题
	Content string `json:"content"` // 内容
	UserID  uint   `json:"userId"` // 用户 ID
}

// ArticleQueryParams 文章查询参数
type ArticleQueryParams struct {
	Page          int  `form:"page"`      // 页码
	PageSize      int  `form:"pageSize"`  // 每页条数
	WithRelations bool `form:"withRelations"` // 是否加载关联
	Title         string `form:"title"` // 标题
	UserID        uint   `form:"userId"` // 用户 ID
}

// ArticleResponse 文章响应
type ArticleResponse struct {
	ID        uint       `json:"id"`         // ID
	CreatedAt string     `json:"createdAt"`  // 创建时间
	UpdatedAt string     `json:"updatedAt"`  // 更新时间
	Title     string     `json:"title"`      // 标题
	Content   string     `json:"content"`    // 内容
	User      *UserResponse `json:"user"`    // 用户
	Comments  []*CommentResponse `json:"comments"` // 评论列表
	Tags      []*TagResponse `json:"tags"`   // 标签列表
}

// ArticleListResponse 文章列表响应
type ArticleListResponse struct {
	Total int64             `json:"total"`  // 总数
	List  []*ArticleResponse `json:"list"`  // 列表
}

// UserResponse 用户响应
type UserResponse struct {
	ID   uint   `json:"id"`   // ID
	Name string `json:"name"` // 名称
}

// CommentResponse 评论响应
type CommentResponse struct {
	ID      uint   `json:"id"`      // ID
	Content string `json:"content"` // 内容
}

// TagResponse 标签响应
type TagResponse struct {
	ID   uint   `json:"id"`   // ID
	Name string `json:"name"` // 名称
}
