package response

import "time"

type OrderResponse struct {
	ID              int                 `json:"id"`
	UserID          int                 `json:"user_id"`
	VoucherCode     *string             `json:"voucher_code"`
	DiscountAmount  float64             `json:"discount_amount"`
	TotalAmount     float64             `json:"total_amount"`
	Status          string              `json:"status"`
	ShippingAddress string              `json:"shipping_address"`
	CreatedAt       time.Time           `json:"created_at"`
	Items           []OrderItemResponse `json:"items"`
	Payment         *PaymentResponse    `json:"payment,omitempty"`
}

type OrderItemResponse struct {
	ID          int                    `json:"id"`
	ProductName string                 `json:"product_name"`
	Variant     ProductVariantResponse `json:"variant"`
	Price       float64                `json:"price"`
	Quantity    int                    `json:"quantity"`
	SubTotal    float64                `json:"sub_total"`
}

type PaymentResponse struct {
	ID            int        `json:"id"`
	PaymentMethod string     `json:"payment_method"`
	PaymentStatus string     `json:"payment_status"`
	PaidAt        *time.Time `json:"paid_at"`
}
