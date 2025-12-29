package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/leehai1107/chophimco-server/pkg/apiwrapper"
	"github.com/leehai1107/chophimco-server/pkg/logger"
	"github.com/leehai1107/chophimco-server/service/chophimco/model/request"
)

type ISellerHandler interface {
	// Seller profile management
	CreateSellerProfile(ctx *gin.Context)
	GetSellerProfile(ctx *gin.Context)
	UpdateSellerProfile(ctx *gin.Context)
	GetSellerByID(ctx *gin.Context)

	// Seller product management
	CreateSellerProduct(ctx *gin.Context)
	GetSellerProducts(ctx *gin.Context)
	UpdateSellerProduct(ctx *gin.Context)
	DeleteSellerProduct(ctx *gin.Context)

	// Product image management
	UploadProductImage(ctx *gin.Context)
	DeleteProductImage(ctx *gin.Context)
	SetPrimaryImage(ctx *gin.Context)

	// Admin - seller verification
	GetPendingSellers(ctx *gin.Context)
	VerifySeller(ctx *gin.Context)
	RejectSeller(ctx *gin.Context)

	// Admin - product approval
	GetPendingProducts(ctx *gin.Context)
	ApproveProduct(ctx *gin.Context)
	RejectProduct(ctx *gin.Context)

	// Seller reviews
	GetSellerReviews(ctx *gin.Context)
	CreateSellerReview(ctx *gin.Context)
}

// CreateSellerProfile godoc
// @Summary Create seller profile
// @Description Create a new seller profile (requires seller role)
// @Tags seller
// @Accept json
// @Produce json
// @Param request body request.CreateSellerProfile true "Seller profile data"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/seller/profile [post]
func (h *Handler) CreateSellerProfile(ctx *gin.Context) {
	var req request.CreateSellerProfile
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	// TODO: Get user ID from JWT token
	userID := ctx.GetInt("user_id")
	req.UserID = userID

	profile, err := h.sellerUsecase.CreateSellerProfile(ctx, req)
	if err != nil {
		logger.EnhanceWith(ctx).Errorw("Failed to create seller profile", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to create seller profile")
		return
	}

	apiwrapper.SendSuccess(ctx, profile)
}

// GetSellerProfile godoc
// @Summary Get seller profile
// @Description Get authenticated seller's profile
// @Tags seller
// @Produce json
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/seller/profile [get]
func (h *Handler) GetSellerProfile(ctx *gin.Context) {
	userID := ctx.GetInt("user_id")

	profile, err := h.sellerUsecase.GetSellerProfileByUserID(ctx, userID)
	if err != nil {
		apiwrapper.SendNotFound(ctx, "Seller profile not found")
		return
	}

	apiwrapper.SendSuccess(ctx, profile)
}

// UpdateSellerProfile godoc
// @Summary Update seller profile
// @Description Update seller's profile information
// @Tags seller
// @Accept json
// @Produce json
// @Param request body request.UpdateSellerProfile true "Update data"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/seller/profile [put]
func (h *Handler) UpdateSellerProfile(ctx *gin.Context) {
	var req request.UpdateSellerProfile
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	userID := ctx.GetInt("user_id")
	req.UserID = userID

	profile, err := h.sellerUsecase.UpdateSellerProfile(ctx, req)
	if err != nil {
		logger.EnhanceWith(ctx).Errorw("Failed to update seller profile", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to update seller profile")
		return
	}

	apiwrapper.SendSuccess(ctx, profile)
}

// GetSellerByID godoc
// @Summary Get seller by ID
// @Description Get public seller information by ID
// @Tags seller
// @Produce json
// @Param id path int true "Seller ID"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/seller/{id} [get]
func (h *Handler) GetSellerByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid seller ID")
		return
	}

	seller, err := h.sellerUsecase.GetSellerByID(ctx, id)
	if err != nil {
		apiwrapper.SendNotFound(ctx, "Seller not found")
		return
	}

	apiwrapper.SendSuccess(ctx, seller)
}

// CreateSellerProduct godoc
// @Summary Create product as seller
// @Description Seller creates a new product (pending approval)
// @Tags seller
// @Accept json
// @Produce json
// @Param request body request.CreateProduct true "Product data"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/seller/product [post]
func (h *Handler) CreateSellerProduct(ctx *gin.Context) {
	var req request.CreateProduct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	userID := ctx.GetInt("user_id")
	req.SellerID = userID

	product, err := h.sellerUsecase.CreateProduct(ctx, req)
	if err != nil {
		logger.EnhanceWith(ctx).Errorw("Failed to create product", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to create product")
		return
	}

	apiwrapper.SendSuccess(ctx, product)
}

// GetSellerProducts godoc
// @Summary Get seller's products
// @Description Get all products belonging to authenticated seller
// @Tags seller
// @Produce json
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/seller/product [get]
func (h *Handler) GetSellerProducts(ctx *gin.Context) {
	userID := ctx.GetInt("user_id")

	products, err := h.sellerUsecase.GetProductsBySeller(ctx, userID)
	if err != nil {
		logger.EnhanceWith(ctx).Errorw("Failed to get products", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to get products")
		return
	}

	apiwrapper.SendSuccess(ctx, products)
}

// UpdateSellerProduct godoc
// @Summary Update seller's product
// @Description Update product details (requires re-approval if major changes)
// @Tags seller
// @Accept json
// @Produce json
// @Param request body request.UpdateProduct true "Update data"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/seller/product [put]
func (h *Handler) UpdateSellerProduct(ctx *gin.Context) {
	var req request.UpdateProduct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	userID := ctx.GetInt("user_id")

	product, err := h.sellerUsecase.UpdateProduct(ctx, userID, req)
	if err != nil {
		logger.EnhanceWith(ctx).Errorw("Failed to update product", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to update product")
		return
	}

	apiwrapper.SendSuccess(ctx, product)
}

// DeleteSellerProduct godoc
// @Summary Delete seller's product
// @Description Delete a product owned by the seller
// @Tags seller
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/seller/product/{id} [delete]
func (h *Handler) DeleteSellerProduct(ctx *gin.Context) {
	productID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid product ID")
		return
	}

	userID := ctx.GetInt("user_id")

	err = h.sellerUsecase.DeleteProduct(ctx, userID, productID)
	if err != nil {
		logger.EnhanceWith(ctx).Errorw("Failed to delete product", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to delete product")
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Product deleted successfully"})
}

// UploadProductImage godoc
// @Summary Upload product image
// @Description Upload an image for a product
// @Tags seller
// @Accept multipart/form-data
// @Produce json
// @Param product_id formData int true "Product ID"
// @Param image formData file true "Image file"
// @Param is_primary formData bool false "Set as primary image"
// @Param display_order formData int false "Display order"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/seller/product/image [post]
func (h *Handler) UploadProductImage(ctx *gin.Context) {
	var req request.UploadProductImage
	if err := ctx.ShouldBind(&req); err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	userID := ctx.GetInt("user_id")

	image, err := h.sellerUsecase.UploadProductImage(ctx, userID, req)
	if err != nil {
		logger.EnhanceWith(ctx).Errorw("Failed to upload image", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to upload image")
		return
	}

	apiwrapper.SendSuccess(ctx, image)
}

// DeleteProductImage godoc
// @Summary Delete product image
// @Description Delete a product image
// @Tags seller
// @Produce json
// @Param id path int true "Image ID"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/seller/product/image/{id} [delete]
func (h *Handler) DeleteProductImage(ctx *gin.Context) {
	imageID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid image ID")
		return
	}

	userID := ctx.GetInt("user_id")

	err = h.sellerUsecase.DeleteProductImage(ctx, userID, imageID)
	if err != nil {
		logger.EnhanceWith(ctx).Errorw("Failed to delete image", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to delete image")
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Image deleted successfully"})
}

// SetPrimaryImage godoc
// @Summary Set primary product image
// @Description Set an image as the primary image for a product
// @Tags seller
// @Accept json
// @Produce json
// @Param request body request.SetPrimaryImage true "Image data"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/seller/product/image/primary [put]
func (h *Handler) SetPrimaryImage(ctx *gin.Context) {
	var req request.SetPrimaryImage
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	userID := ctx.GetInt("user_id")

	err := h.sellerUsecase.SetPrimaryImage(ctx, userID, req)
	if err != nil {
		logger.EnhanceWith(ctx).Errorw("Failed to set primary image", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to set primary image")
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Primary image updated"})
}

// GetPendingSellers godoc
// @Summary Get pending seller verifications
// @Description Admin - Get all sellers pending verification
// @Tags admin
// @Produce json
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/admin/seller/pending [get]
func (h *Handler) GetPendingSellers(ctx *gin.Context) {
	sellers, err := h.sellerUsecase.GetPendingSellers(ctx)
	if err != nil {
		logger.EnhanceWith(ctx).Errorw("Failed to get pending sellers", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to get pending sellers")
		return
	}

	apiwrapper.SendSuccess(ctx, sellers)
}

// VerifySeller godoc
// @Summary Verify seller
// @Description Admin - Verify a seller profile
// @Tags admin
// @Accept json
// @Produce json
// @Param request body request.VerifySeller true "Verification data"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/admin/seller/verify [post]
func (h *Handler) VerifySeller(ctx *gin.Context) {
	var req request.VerifySeller
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	err := h.sellerUsecase.VerifySeller(ctx, req)
	if err != nil {
		logger.EnhanceWith(ctx).Errorw("Failed to verify seller", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to verify seller")
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Seller verified successfully"})
}

// RejectSeller godoc
// @Summary Reject seller
// @Description Admin - Reject a seller profile
// @Tags admin
// @Accept json
// @Produce json
// @Param request body request.RejectSeller true "Rejection data"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/admin/seller/reject [post]
func (h *Handler) RejectSeller(ctx *gin.Context) {
	var req request.RejectSeller
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	err := h.sellerUsecase.RejectSeller(ctx, req)
	if err != nil {
		logger.EnhanceWith(ctx).Errorw("Failed to reject seller", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to reject seller")
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Seller rejected"})
}

// GetPendingProducts godoc
// @Summary Get pending products
// @Description Admin - Get all products pending approval
// @Tags admin
// @Produce json
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/admin/product/pending [get]
func (h *Handler) GetPendingProducts(ctx *gin.Context) {
	products, err := h.sellerUsecase.GetPendingProducts(ctx)
	if err != nil {
		logger.EnhanceWith(ctx).Errorw("Failed to get pending products", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to get pending products")
		return
	}

	apiwrapper.SendSuccess(ctx, products)
}

// ApproveProduct godoc
// @Summary Approve product
// @Description Admin - Approve a product for listing
// @Tags admin
// @Accept json
// @Produce json
// @Param request body request.ApproveProduct true "Approval data"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/admin/product/approve [post]
func (h *Handler) ApproveProduct(ctx *gin.Context) {
	var req request.ApproveProduct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	err := h.sellerUsecase.ApproveProduct(ctx, req)
	if err != nil {
		logger.EnhanceWith(ctx).Errorw("Failed to approve product", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to approve product")
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Product approved successfully"})
}

// RejectProduct godoc
// @Summary Reject product
// @Description Admin - Reject a product with reason
// @Tags admin
// @Accept json
// @Produce json
// @Param request body request.RejectProduct true "Rejection data"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/admin/product/reject [post]
func (h *Handler) RejectProduct(ctx *gin.Context) {
	var req request.RejectProduct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	err := h.sellerUsecase.RejectProduct(ctx, req)
	if err != nil {
		logger.EnhanceWith(ctx).Errorw("Failed to reject product", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to reject product")
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Product rejected"})
}

// GetSellerReviews godoc
// @Summary Get seller reviews
// @Description Get all reviews for a seller
// @Tags seller
// @Produce json
// @Param seller_id query int true "Seller ID"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/seller/reviews [get]
func (h *Handler) GetSellerReviews(ctx *gin.Context) {
	sellerID, err := strconv.Atoi(ctx.Query("seller_id"))
	if err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid seller ID")
		return
	}

	reviews, err := h.sellerUsecase.GetSellerReviews(ctx, sellerID)
	if err != nil {
		logger.EnhanceWith(ctx).Errorw("Failed to get seller reviews", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to get seller reviews")
		return
	}

	apiwrapper.SendSuccess(ctx, reviews)
}

// CreateSellerReview godoc
// @Summary Create seller review
// @Description Create a review for a seller after order completion
// @Tags seller
// @Accept json
// @Produce json
// @Param request body request.CreateSellerReview true "Review data"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/seller/reviews [post]
func (h *Handler) CreateSellerReview(ctx *gin.Context) {
	var req request.CreateSellerReview
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	userID := ctx.GetInt("user_id")
	req.BuyerID = userID

	review, err := h.sellerUsecase.CreateSellerReview(ctx, req)
	if err != nil {
		logger.EnhanceWith(ctx).Errorw("Failed to create seller review", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to create seller review")
		return
	}

	apiwrapper.SendSuccess(ctx, review)
}
