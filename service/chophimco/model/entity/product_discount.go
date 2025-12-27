package entity

import (
	"time"
)

type ProductDiscount struct {
	ID            int        `gorm:"primaryKey;column:id;autoIncrement"`
	ProductID     int        `gorm:"column:product_id;not null"`
	DiscountType  string     `gorm:"column:discount_type;not null"` // percent | fixed
	DiscountValue float64    `gorm:"column:discount_value;not null"`
	StartAt       *time.Time `gorm:"column:start_at"`
	EndAt         *time.Time `gorm:"column:end_at"`
	IsActive      bool       `gorm:"column:is_active;default:true"`
	CreatedAt     time.Time  `gorm:"column:created_at;default:now()"`

	// Relations
	Product *Product `gorm:"foreignKey:ProductID;references:ID"`
}
