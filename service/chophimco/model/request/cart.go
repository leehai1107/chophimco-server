package request

type AddToCart struct {
	ProductVariantID int `json:"product_variant_id" binding:"required"`
	Quantity         int `json:"quantity" binding:"required,gt=0"`
}

type UpdateCartItem struct {
	CartItemID int `json:"cart_item_id" binding:"required"`
	Quantity   int `json:"quantity" binding:"required,gt=0"`
}
