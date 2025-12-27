package entity

import (
	"time"
)

type Order struct {
	ID              int       `gorm:"primaryKey;column:id;autoIncrement"`
	UserID          int       `gorm:"column:user_id;not null"`
	VoucherID       *int      `gorm:"column:voucher_id"`
	DiscountAmount  float64   `gorm:"column:discount_amount;default:0"`
	TotalAmount     float64   `gorm:"column:total_amount;not null"`
	Status          string    `gorm:"column:status;not null"` // pending, paid, shipped, completed, cancelled
	ShippingAddress string    `gorm:"column:shipping_address;type:text"`
	CreatedAt       time.Time `gorm:"column:created_at;default:now()"`

	// Relations
	User       *User       `gorm:"foreignKey:UserID;references:ID"`
	Voucher    *Voucher    `gorm:"foreignKey:VoucherID;references:ID"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID"`
	Payment    *Payment    `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID               int     `gorm:"primaryKey;column:id;autoIncrement"`
	OrderID          int     `gorm:"column:order_id;not null"`
	ProductVariantID int     `gorm:"column:product_variant_id;not null"`
	Price            float64 `gorm:"column:price;not null"`
	Quantity         int     `gorm:"column:quantity;not null;check:quantity > 0"`

	// Relations
	Order          *Order          `gorm:"foreignKey:OrderID;references:ID"`
	ProductVariant *ProductVariant `gorm:"foreignKey:ProductVariantID;references:ID"`
}
