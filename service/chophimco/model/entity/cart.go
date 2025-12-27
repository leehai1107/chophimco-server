package entity

import (
	"time"
)

type Cart struct {
	ID        int       `gorm:"primaryKey;column:id;autoIncrement"`
	UserID    int       `gorm:"column:user_id;not null"`
	CreatedAt time.Time `gorm:"column:created_at;default:now()"`

	// Relations
	User      *User      `gorm:"foreignKey:UserID;references:ID"`
	CartItems []CartItem `gorm:"foreignKey:CartID"`
}

type CartItem struct {
	ID               int `gorm:"primaryKey;column:id;autoIncrement"`
	CartID           int `gorm:"column:cart_id;not null"`
	ProductVariantID int `gorm:"column:product_variant_id;not null"`
	Quantity         int `gorm:"column:quantity;not null;check:quantity > 0"`

	// Relations
	Cart           *Cart           `gorm:"foreignKey:CartID;references:ID"`
	ProductVariant *ProductVariant `gorm:"foreignKey:ProductVariantID;references:ID"`
}
