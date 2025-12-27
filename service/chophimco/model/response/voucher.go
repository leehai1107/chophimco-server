package response

import "time"

type VoucherResponse struct {
	ID               int        `json:"id"`
	Code             string     `json:"code"`
	Description      string     `json:"description"`
	DiscountType     string     `json:"discount_type"`
	DiscountValue    float64    `json:"discount_value"`
	MinOrderValue    float64    `json:"min_order_value"`
	MaxDiscountValue *float64   `json:"max_discount_value"`
	UsageLimit       *int       `json:"usage_limit"`
	UsagePerUser     int        `json:"usage_per_user"`
	UsedCount        int        `json:"used_count"`
	StartAt          *time.Time `json:"start_at"`
	EndAt            *time.Time `json:"end_at"`
	IsActive         bool       `json:"is_active"`
	CreatedAt        time.Time  `json:"created_at"`
}
