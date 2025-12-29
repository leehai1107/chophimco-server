package request

type CreateSellerProfile struct {
	UserID          int    `json:"user_id"`
	ShopName        string `json:"shop_name" binding:"required"`
	ShopDescription string `json:"shop_description"`
	BusinessAddress string `json:"business_address"`
	BusinessPhone   string `json:"business_phone"`
	LogoURL         string `json:"logo_url"`
}

type UpdateSellerProfile struct {
	UserID          int    `json:"user_id"`
	ShopName        string `json:"shop_name"`
	ShopDescription string `json:"shop_description"`
	BusinessAddress string `json:"business_address"`
	BusinessPhone   string `json:"business_phone"`
	LogoURL         string `json:"logo_url"`
}

type UploadProductImage struct {
	ProductID    int    `form:"product_id" binding:"required"`
	ImageURL     string `form:"image_url" binding:"required"`
	IsPrimary    bool   `form:"is_primary"`
	DisplayOrder int    `form:"display_order"`
	AltText      string `form:"alt_text"`
}

type SetPrimaryImage struct {
	ProductID int `json:"product_id" binding:"required"`
	ImageID   int `json:"image_id" binding:"required"`
}

type VerifySeller struct {
	UserID int `json:"user_id" binding:"required"`
}

type RejectSeller struct {
	UserID int    `json:"user_id" binding:"required"`
	Reason string `json:"reason"`
}

type ApproveProduct struct {
	ProductID int `json:"product_id" binding:"required"`
}

type RejectProduct struct {
	ProductID int    `json:"product_id" binding:"required"`
	Reason    string `json:"reason" binding:"required"`
}

type CreateSellerReview struct {
	BuyerID  int    `json:"buyer_id"`
	SellerID int    `json:"seller_id" binding:"required"`
	OrderID  int    `json:"order_id" binding:"required"`
	Rating   int    `json:"rating" binding:"required,min=1,max=5"`
	Comment  string `json:"comment"`
}
