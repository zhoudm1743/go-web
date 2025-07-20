package models

import (
	"time"
	"github.com/zhoudm1743/go-web/core/facades"
	"gorm.io/gorm"
)

// Article 文章
type Article struct {
	ID        uint           `gorm:"primarykey" json:"id"`                // 主键ID
	CreatedAt time.Time      `json:"createdAt"`                           // 创建时间
	UpdatedAt time.Time      `json:"updatedAt"`                           // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                      // 删除时间
	Title     string         `gorm:"column:title" json:"title"`           // 标题
	Content   string         `gorm:"column:content" json:"content"`       // 内容
	UserID    uint           `gorm:"column:user_id" json:"userId"`        // User外键
	User      User           `gorm:"foreignKey:UserID;references:ID" json:"user"`  // 用户
	Comments  []Comment      `gorm:"foreignKey:ArticleID" json:"comments"`        // 评论列表
	Tags      []Tag          `gorm:"many2many:article_tags;foreignKey:article_id;references:tag_id" json:"tags"` // 标签列表
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
