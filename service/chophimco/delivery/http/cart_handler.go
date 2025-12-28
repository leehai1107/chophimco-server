package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/leehai1107/chophimco-server/pkg/apiwrapper"
	"github.com/leehai1107/chophimco-server/service/chophimco/model/request"
)

type ICartHandler interface {
	GetCart(ctx *gin.Context)
	AddToCart(ctx *gin.Context)
	UpdateCartItem(ctx *gin.Context)
	RemoveFromCart(ctx *gin.Context)
	ClearCart(ctx *gin.Context)
}

// GetCart godoc
// @Summary Get cart
// @Description Get user's shopping cart
// @Tags cart
// @Produce json
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/cart [get]
func (h *Handler) GetCart(ctx *gin.Context) {
	userIDStr := ctx.GetString("user_id")
	userID, _ := strconv.Atoi(userIDStr)

	cart, err := h.cartUsecase.GetCart(ctx, userID)
	if err != nil {
		apiwrapper.SendInternalError(ctx, "Failed to get cart")
		return
	}

	apiwrapper.SendSuccess(ctx, cart)
}

// AddToCart godoc
// @Summary Add to cart
// @Description Add product to cart
// @Tags cart
// @Accept json
// @Produce json
// @Param request body request.AddToCart true "Cart item"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/cart/add [post]
func (h *Handler) AddToCart(ctx *gin.Context) {
	userIDStr := ctx.GetString("user_id")
	userID, _ := strconv.Atoi(userIDStr)

	var req request.AddToCart
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	if err := h.cartUsecase.AddToCart(ctx, userID, req); err != nil {
		apiwrapper.SendInternalError(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Added to cart"})
}

// UpdateCartItem godoc
// @Summary Update cart item
// @Description Update cart item quantity
// @Tags cart
// @Accept json
// @Produce json
// @Param request body request.UpdateCartItem true "Cart item update"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/cart/update [put]
func (h *Handler) UpdateCartItem(ctx *gin.Context) {
	var req request.UpdateCartItem
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	if err := h.cartUsecase.UpdateCartItem(ctx, req); err != nil {
		apiwrapper.SendInternalError(ctx, "Failed to update cart item")
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Cart item updated"})
}

// RemoveFromCart godoc
// @Summary Remove from cart
// @Description Remove item from cart
// @Tags cart
// @Param id path int true "Cart Item ID"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/cart/{id} [delete]
func (h *Handler) RemoveFromCart(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid cart item ID")
		return
	}

	if err := h.cartUsecase.RemoveFromCart(ctx, id); err != nil {
		apiwrapper.SendInternalError(ctx, "Failed to remove item")
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Item removed from cart"})
}

// ClearCart godoc
// @Summary Clear cart
// @Description Clear all items from cart
// @Tags cart
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/cart/clear [delete]
func (h *Handler) ClearCart(ctx *gin.Context) {
	userIDStr := ctx.GetString("user_id")
	userID, _ := strconv.Atoi(userIDStr)

	if err := h.cartUsecase.ClearCart(ctx, userID); err != nil {
		apiwrapper.SendInternalError(ctx, "Failed to clear cart")
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Cart cleared"})
}
