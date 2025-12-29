package entity

import (
	"time"
)

type Product struct {
	ID              int        `gorm:"primaryKey;column:id;autoIncrement"`
	SellerID        int        `gorm:"column:seller_id;not null"`
	Name            string     `gorm:"column:name;not null"`
	CategoryID      *int       `gorm:"column:category_id"`
	BrandID         *int       `gorm:"column:brand_id"`
	Description     string     `gorm:"column:description;type:text"`
	BasePrice       float64    `gorm:"column:base_price;not null"`
	ApprovalStatus  string     `gorm:"column:approval_status;default:pending"`
	RejectionReason string     `gorm:"column:rejection_reason;type:text"`
	IsActive        bool       `gorm:"column:is_active;default:true"`
	CreatedAt       time.Time  `gorm:"column:created_at;default:now()"`
	ApprovedAt      *time.Time `gorm:"column:approved_at"`

	// Relations
	Seller   *User            `gorm:"foreignKey:SellerID;references:ID"`
	Category *Category        `gorm:"foreignKey:CategoryID;references:ID"`
	Brand    *Brand           `gorm:"foreignKey:BrandID;references:ID"`
	Variants []ProductVariant `gorm:"foreignKey:ProductID"`
	Images   []ProductImage   `gorm:"foreignKey:ProductID"`
}
