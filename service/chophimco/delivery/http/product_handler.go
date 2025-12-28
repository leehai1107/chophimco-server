package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/leehai1107/chophimco-server/pkg/apiwrapper"
	"github.com/leehai1107/chophimco-server/pkg/logger"
	"github.com/leehai1107/chophimco-server/service/chophimco/model/request"
)

type IProductHandler interface {
	GetAllProducts(ctx *gin.Context)
	GetProductByID(ctx *gin.Context)
	GetProductsByCategory(ctx *gin.Context)
	GetProductsByBrand(ctx *gin.Context)
	CreateProduct(ctx *gin.Context)
	UpdateProduct(ctx *gin.Context)
	DeleteProduct(ctx *gin.Context)
	CreateProductVariant(ctx *gin.Context)
	UpdateProductVariant(ctx *gin.Context)
}

// GetAllProducts godoc
// @Summary Get all products
// @Description Get all active products
// @Tags product
// @Produce json
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/product/all [get]
func (h *Handler) GetAllProducts(ctx *gin.Context) {
	products, err := h.productUsecase.GetAllProducts(ctx)
	if err != nil {
		logger.EnhanceWith(ctx).Errorw("Failed to get products", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to get products")
		return
	}

	apiwrapper.SendSuccess(ctx, products)
}

// GetProductByID godoc
// @Summary Get product by ID
// @Description Get product details by ID
// @Tags product
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/product/{id} [get]
func (h *Handler) GetProductByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid product ID")
		return
	}

	product, err := h.productUsecase.GetProductByID(ctx, id)
	if err != nil {
		apiwrapper.SendNotFound(ctx, "Product not found")
		return
	}

	apiwrapper.SendSuccess(ctx, product)
}

// GetProductsByCategory godoc
// @Summary Get products by category
// @Description Get all products in a category
// @Tags product
// @Produce json
// @Param category_id query int true "Category ID"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/product/category [get]
func (h *Handler) GetProductsByCategory(ctx *gin.Context) {
	categoryID, err := strconv.Atoi(ctx.Query("category_id"))
	if err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid category ID")
		return
	}

	products, err := h.productUsecase.GetProductsByCategory(ctx, categoryID)
	if err != nil {
		apiwrapper.SendInternalError(ctx, "Failed to get products")
		return
	}

	apiwrapper.SendSuccess(ctx, products)
}

// GetProductsByBrand godoc
// @Summary Get products by brand
// @Description Get all products of a brand
// @Tags product
// @Produce json
// @Param brand_id query int true "Brand ID"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/product/brand [get]
func (h *Handler) GetProductsByBrand(ctx *gin.Context) {
	brandID, err := strconv.Atoi(ctx.Query("brand_id"))
	if err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid brand ID")
		return
	}

	products, err := h.productUsecase.GetProductsByBrand(ctx, brandID)
	if err != nil {
		apiwrapper.SendInternalError(ctx, "Failed to get products")
		return
	}

	apiwrapper.SendSuccess(ctx, products)
}

// CreateProduct godoc
// @Summary Create new product
// @Description Create a new product (Admin only)
// @Tags product
// @Accept json
// @Produce json
// @Param request body request.CreateProduct true "Product information"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/product/create [post]
func (h *Handler) CreateProduct(ctx *gin.Context) {
	var req request.CreateProduct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	if err := h.productUsecase.CreateProduct(ctx, req); err != nil {
		apiwrapper.SendInternalError(ctx, "Failed to create product")
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Product created successfully"})
}

// UpdateProduct godoc
// @Summary Update product
// @Description Update product information (Admin only)
// @Tags product
// @Accept json
// @Produce json
// @Param request body request.UpdateProduct true "Product information"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/product/update [put]
func (h *Handler) UpdateProduct(ctx *gin.Context) {
	var req request.UpdateProduct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	if err := h.productUsecase.UpdateProduct(ctx, req); err != nil {
		apiwrapper.SendInternalError(ctx, "Failed to update product")
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Product updated successfully"})
}

// DeleteProduct godoc
// @Summary Delete product
// @Description Soft delete a product (Admin only)
// @Tags product
// @Param id path int true "Product ID"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/product/{id} [delete]
func (h *Handler) DeleteProduct(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid product ID")
		return
	}

	if err := h.productUsecase.DeleteProduct(ctx, id); err != nil {
		apiwrapper.SendInternalError(ctx, "Failed to delete product")
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Product deleted successfully"})
}

// CreateProductVariant godoc
// @Summary Create product variant
// @Description Create a new product variant (Admin only)
// @Tags product
// @Accept json
// @Produce json
// @Param request body request.CreateProductVariant true "Variant information"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/product/variant/create [post]
func (h *Handler) CreateProductVariant(ctx *gin.Context) {
	var req request.CreateProductVariant
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	if err := h.productUsecase.CreateProductVariant(ctx, req); err != nil {
		if err.Error() == "SKU already exists" {
			apiwrapper.SendBadRequest(ctx, err.Error())
			return
		}
		apiwrapper.SendInternalError(ctx, "Failed to create variant")
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Variant created successfully"})
}

// UpdateProductVariant godoc
// @Summary Update product variant
// @Description Update product variant information (Admin only)
// @Tags product
// @Accept json
// @Produce json
// @Param request body request.UpdateProductVariant true "Variant information"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/product/variant/update [put]
func (h *Handler) UpdateProductVariant(ctx *gin.Context) {
	var req request.UpdateProductVariant
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	if err := h.productUsecase.UpdateProductVariant(ctx, req); err != nil {
		apiwrapper.SendInternalError(ctx, "Failed to update variant")
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Variant updated successfully"})
}
