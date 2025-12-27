package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/leehai1107/chophimco-server/pkg/apiwrapper"
	"github.com/leehai1107/chophimco-server/service/chophimco/model/request"
)

type IOrderHandler interface {
	CreateOrder(ctx *gin.Context)
	GetOrderByID(ctx *gin.Context)
	GetMyOrders(ctx *gin.Context)
	UpdateOrderStatus(ctx *gin.Context)
}

// CreateOrder godoc
// @Summary Create order
// @Description Create a new order from cart
// @Tags order
// @Accept json
// @Produce json
// @Param request body request.CreateOrder true "Order information"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /chophimco/api/v1/order/create [post]
func (h *Handler) CreateOrder(ctx *gin.Context) {
	userIDStr := ctx.GetString("user_id")
	userID, _ := strconv.Atoi(userIDStr)

	var req request.CreateOrder
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	order, err := h.orderUsecase.CreateOrder(ctx, userID, req)
	if err != nil {
		apiwrapper.SendInternalError(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, order)
}

// GetOrderByID godoc
// @Summary Get order by ID
// @Description Get order details by ID
// @Tags order
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /chophimco/api/v1/order/{id} [get]
func (h *Handler) GetOrderByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid order ID")
		return
	}

	order, err := h.orderUsecase.GetOrderByID(ctx, id)
	if err != nil {
		apiwrapper.SendNotFound(ctx, "Order not found")
		return
	}

	apiwrapper.SendSuccess(ctx, order)
}

// GetMyOrders godoc
// @Summary Get user orders
// @Description Get all orders of current user
// @Tags order
// @Produce json
// @Success 200 {object} apiwrapper.APIResponse
// @Router /chophimco/api/v1/order/my-orders [get]
func (h *Handler) GetMyOrders(ctx *gin.Context) {
	userIDStr := ctx.GetString("user_id")
	userID, _ := strconv.Atoi(userIDStr)

	orders, err := h.orderUsecase.GetUserOrders(ctx, userID)
	if err != nil {
		apiwrapper.SendInternalError(ctx, "Failed to get orders")
		return
	}

	apiwrapper.SendSuccess(ctx, orders)
}

// UpdateOrderStatus godoc
// @Summary Update order status
// @Description Update order status (Admin only)
// @Tags order
// @Accept json
// @Produce json
// @Param request body request.UpdateOrderStatus true "Order status"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /chophimco/api/v1/order/status [put]
func (h *Handler) UpdateOrderStatus(ctx *gin.Context) {
	var req request.UpdateOrderStatus
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	if err := h.orderUsecase.UpdateOrderStatus(ctx, req); err != nil {
		apiwrapper.SendInternalError(ctx, "Failed to update order status")
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Order status updated"})
}
