package response

type CartResponse struct {
	ID        int                `json:"id"`
	UserID    int                `json:"user_id"`
	Items     []CartItemResponse `json:"items"`
	TotalItem int                `json:"total_items"`
	SubTotal  float64            `json:"sub_total"`
}

type CartItemResponse struct {
	ID          int                    `json:"id"`
	ProductName string                 `json:"product_name"`
	Variant     ProductVariantResponse `json:"variant"`
	Quantity    int                    `json:"quantity"`
	Price       float64                `json:"price"`
	SubTotal    float64                `json:"sub_total"`
}
