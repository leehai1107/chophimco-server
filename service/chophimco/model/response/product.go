package response

import "time"

type ProductResponse struct {
	ID          int                      `json:"id"`
	Name        string                   `json:"name"`
	Category    *string                  `json:"category"`
	Brand       *string                  `json:"brand"`
	Description string                   `json:"description"`
	BasePrice   float64                  `json:"base_price"`
	IsActive    bool                     `json:"is_active"`
	CreatedAt   time.Time                `json:"created_at"`
	Variants    []ProductVariantResponse `json:"variants,omitempty"`
}

type ProductVariantResponse struct {
	ID             int     `json:"id"`
	ProductID      int     `json:"product_id"`
	Switch         *string `json:"switch"`
	Layout         string  `json:"layout"`
	ConnectionType string  `json:"connection_type"`
	Hotswap        bool    `json:"hotswap"`
	LedType        string  `json:"led_type"`
	Price          float64 `json:"price"`
	Stock          int     `json:"stock"`
	SKU            string  `json:"sku"`
}
