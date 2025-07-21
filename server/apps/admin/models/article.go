package models

import (
	"time"

	"gorm.io/gorm"
)

// Article 文章
type Article struct {
	ID        uint           `gorm:"primarykey" json:"id"` // 主键ID
	CreatedAt time.Time      `json:"createdAt"`            // 创建时间
	UpdatedAt time.Time      `json:"updatedAt"`            // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`       // 删除时间

	Title string `gorm:"column:title" json:"title"` // 标题

	Category *Category `gorm:"foreignKey:CategoryID;references:ID" json:"category"` // 分类

	CategoryID uint `gorm:"column:category_id" json:"categoryID"` // Category外键

	Author string `gorm:"column:author" json:"author"` // 作者

}

// TableName 设置表名
func (Article) TableName() string {
	return "articles"
}

// LoadRelations 关系预加载
func (m *Article) LoadRelations(db *gorm.DB) *gorm.DB {
	query := db

	query = query.Preload("Category")

	return query
}
