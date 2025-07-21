package models

import (
	"time"
	"gorm.io/gorm"
)

// Product 产品
type Product struct {
	ID        uint           `gorm:"primarykey" json:"id"`                // 主键ID
	CreatedAt time.Time      `json:"createdAt"`                           // 创建时间
	UpdatedAt time.Time      `json:"updatedAt"`                           // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                      // 删除时间

	CategoryID string `gorm:"column:categoryID" json:"categoryID"` // 产品分类

	Name string `gorm:"column:name" json:"name"` // 名称

	Sort int64 `gorm:"column:sort" json:"sort"` // 排序

}

// TableName 设置表名
func (Product) TableName() string {
	return "products"
}


