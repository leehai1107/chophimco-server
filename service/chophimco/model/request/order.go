package request

type CreateOrder struct {
	VoucherCode     string `json:"voucher_code"`
	ShippingAddress string `json:"shipping_address" binding:"required"`
}

type UpdateOrderStatus struct {
	OrderID int    `json:"order_id" binding:"required"`
	Status  string `json:"status" binding:"required,oneof=pending paid shipped completed cancelled"`
}
