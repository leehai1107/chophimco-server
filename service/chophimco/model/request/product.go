package request

type CreateProduct struct {
	SellerID    int     `json:"seller_id"`
	Name        string  `json:"name" binding:"required"`
	CategoryID  *int    `json:"category_id"`
	BrandID     *int    `json:"brand_id"`
	Description string  `json:"description"`
	BasePrice   float64 `json:"base_price" binding:"required,gt=0"`
}

type UpdateProduct struct {
	ID          int     `json:"id" binding:"required"`
	Name        string  `json:"name"`
	CategoryID  *int    `json:"category_id"`
	BrandID     *int    `json:"brand_id"`
	Description string  `json:"description"`
	BasePrice   float64 `json:"base_price" binding:"gt=0"`
	IsActive    *bool   `json:"is_active"`
}

type CreateProductVariant struct {
	ProductID      int     `json:"product_id" binding:"required"`
	SwitchID       *int    `json:"switch_id"`
	Layout         string  `json:"layout"`
	ConnectionType string  `json:"connection_type"`
	Hotswap        bool    `json:"hotswap"`
	LedType        string  `json:"led_type"`
	Price          float64 `json:"price" binding:"required,gt=0"`
	Stock          int     `json:"stock" binding:"gte=0"`
	SKU            string  `json:"sku" binding:"required"`
}

type UpdateProductVariant struct {
	ID             int     `json:"id" binding:"required"`
	SwitchID       *int    `json:"switch_id"`
	Layout         string  `json:"layout"`
	ConnectionType string  `json:"connection_type"`
	Hotswap        *bool   `json:"hotswap"`
	LedType        string  `json:"led_type"`
	Price          float64 `json:"price" binding:"gt=0"`
	Stock          *int    `json:"stock" binding:"gte=0"`
}
