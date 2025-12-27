package entity

import (
	"time"
)

type Voucher struct {
	ID               int        `gorm:"primaryKey;column:id;autoIncrement"`
	Code             string     `gorm:"column:code;unique;not null"`
	Description      string     `gorm:"column:description;type:text"`
	DiscountType     string     `gorm:"column:discount_type;not null"` // percent | fixed
	DiscountValue    float64    `gorm:"column:discount_value;not null"`
	MinOrderValue    float64    `gorm:"column:min_order_value;default:0"`
	MaxDiscountValue *float64   `gorm:"column:max_discount_value"`
	UsageLimit       *int       `gorm:"column:usage_limit"`
	UsagePerUser     int        `gorm:"column:usage_per_user;default:1"`
	UsedCount        int        `gorm:"column:used_count;default:0"`
	StartAt          *time.Time `gorm:"column:start_at"`
	EndAt            *time.Time `gorm:"column:end_at"`
	IsActive         bool       `gorm:"column:is_active;default:true"`
	CreatedAt        time.Time  `gorm:"column:created_at;default:now()"`
}

type UserVoucher struct {
	ID        int `gorm:"primaryKey;column:id;autoIncrement"`
	UserID    int `gorm:"column:user_id;not null"`
	VoucherID int `gorm:"column:voucher_id;not null"`
	UsedCount int `gorm:"column:used_count;default:0"`

	// Relations
	User    *User    `gorm:"foreignKey:UserID;references:ID"`
	Voucher *Voucher `gorm:"foreignKey:VoucherID;references:ID"`
}
