package models

import (
	"time"

	"gorm.io/gorm"
)

// Category 分类
type Category struct {
	ID        uint           `gorm:"primarykey" json:"id"` // 主键ID
	CreatedAt time.Time      `json:"createdAt"`            // 创建时间
	UpdatedAt time.Time      `json:"updatedAt"`            // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`       // 删除时间

	Name string `gorm:"column:name" json:"name"` // 分类名称
}

// TableName 设置表名
func (Category) TableName() string {
	return "categories"
}
