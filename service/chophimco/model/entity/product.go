package entity

import (
	"time"
)

type Product struct {
	ID          int       `gorm:"primaryKey;column:id;autoIncrement"`
	Name        string    `gorm:"column:name;not null"`
	CategoryID  *int      `gorm:"column:category_id"`
	BrandID     *int      `gorm:"column:brand_id"`
	Description string    `gorm:"column:description;type:text"`
	BasePrice   float64   `gorm:"column:base_price;not null"`
	IsActive    bool      `gorm:"column:is_active;default:true"`
	CreatedAt   time.Time `gorm:"column:created_at;default:now()"`

	// Relations
	Category *Category        `gorm:"foreignKey:CategoryID;references:ID"`
	Brand    *Brand           `gorm:"foreignKey:BrandID;references:ID"`
	Variants []ProductVariant `gorm:"foreignKey:ProductID"`
}
