package models

import (
	"time"
	"github.com/zhoudm1743/go-web/core/facades"
	"gorm.io/gorm"
)

// Product 产品
type Product struct {
	ID        uint           `gorm:"primarykey" json:"id"`                // 主键ID
	CreatedAt time.Time      `json:"createdAt"`                           // 创建时间
	UpdatedAt time.Time      `json:"updatedAt"`                           // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                      // 删除时间

	ID uint `gorm:"primaryKey;column:id" json:"iD"` // 主键ID

	Name string `gorm:"column:name" json:"name"` // 名称

	Status int `gorm:"column:status" json:"status"` // 状态

	Price float64 `gorm:"column:price" json:"price"` // 价格

}

// TableName 设置表名
func (Product) TableName() string {
	return "products"
}


