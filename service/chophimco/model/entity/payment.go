package entity

import (
	"time"
)

type Payment struct {
	ID            int        `gorm:"primaryKey;column:id;autoIncrement"`
	OrderID       int        `gorm:"column:order_id;not null"`
	PaymentMethod string     `gorm:"column:payment_method"` // COD, Momo, VNPay
	PaymentStatus string     `gorm:"column:payment_status"` // pending, success, failed
	PaidAt        *time.Time `gorm:"column:paid_at"`

	// Relations
	Order *Order `gorm:"foreignKey:OrderID;references:ID"`
}
