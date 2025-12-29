package entity

import (
	"time"
)

type ProductImage struct {
	ID           int       `gorm:"primaryKey;column:id;autoIncrement"`
	ProductID    int       `gorm:"column:product_id;not null"`
	ImageURL     string    `gorm:"column:image_url;type:text;not null"`
	IsPrimary    bool      `gorm:"column:is_primary;default:false"`
	DisplayOrder int       `gorm:"column:display_order;default:0"`
	AltText      string    `gorm:"column:alt_text"`
	CreatedAt    time.Time `gorm:"column:created_at;default:now()"`

	// Relations
	Product *Product `gorm:"foreignKey:ProductID;references:ID"`
}
