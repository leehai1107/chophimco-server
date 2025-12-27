package request

import "time"

type CreateVoucher struct {
	Code             string     `json:"code" binding:"required"`
	Description      string     `json:"description"`
	DiscountType     string     `json:"discount_type" binding:"required,oneof=percent fixed"`
	DiscountValue    float64    `json:"discount_value" binding:"required,gt=0"`
	MinOrderValue    float64    `json:"min_order_value"`
	MaxDiscountValue *float64   `json:"max_discount_value"`
	UsageLimit       *int       `json:"usage_limit"`
	UsagePerUser     int        `json:"usage_per_user"`
	StartAt          *time.Time `json:"start_at"`
	EndAt            *time.Time `json:"end_at"`
}

type UpdateVoucher struct {
	ID               int        `json:"id" binding:"required"`
	Description      string     `json:"description"`
	DiscountType     string     `json:"discount_type" binding:"oneof=percent fixed"`
	DiscountValue    float64    `json:"discount_value" binding:"gt=0"`
	MinOrderValue    float64    `json:"min_order_value"`
	MaxDiscountValue *float64   `json:"max_discount_value"`
	UsageLimit       *int       `json:"usage_limit"`
	UsagePerUser     *int       `json:"usage_per_user"`
	StartAt          *time.Time `json:"start_at"`
	EndAt            *time.Time `json:"end_at"`
	IsActive         *bool      `json:"is_active"`
}

type ApplyVoucher struct {
	Code string `json:"code" binding:"required"`
}
