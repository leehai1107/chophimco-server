package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/leehai1107/chophimco-server/pkg/apiwrapper"
	"github.com/leehai1107/chophimco-server/service/chophimco/model/request"
)

type IVoucherHandler interface {
	GetAllVouchers(ctx *gin.Context)
	GetActiveVouchers(ctx *gin.Context)
	GetVoucherByCode(ctx *gin.Context)
	CreateVoucher(ctx *gin.Context)
	UpdateVoucher(ctx *gin.Context)
	DeleteVoucher(ctx *gin.Context)
	ValidateVoucher(ctx *gin.Context)
}

// GetAllVouchers godoc
// @Summary Get all vouchers
// @Description Get all vouchers (Admin only)
// @Tags voucher
// @Produce json
// @Success 200 {object} apiwrapper.APIResponse
// @Router /chophimco/api/v1/voucher/all [get]
func (h *Handler) GetAllVouchers(ctx *gin.Context) {
	vouchers, err := h.voucherUsecase.GetAllVouchers(ctx)
	if err != nil {
		apiwrapper.SendInternalError(ctx, "Failed to get vouchers")
		return
	}

	apiwrapper.SendSuccess(ctx, vouchers)
}

// GetActiveVouchers godoc
// @Summary Get active vouchers
// @Description Get all active vouchers
// @Tags voucher
// @Produce json
// @Success 200 {object} apiwrapper.APIResponse
// @Router /chophimco/api/v1/voucher/active [get]
func (h *Handler) GetActiveVouchers(ctx *gin.Context) {
	vouchers, err := h.voucherUsecase.GetActiveVouchers(ctx)
	if err != nil {
		apiwrapper.SendInternalError(ctx, "Failed to get vouchers")
		return
	}

	apiwrapper.SendSuccess(ctx, vouchers)
}

// GetVoucherByCode godoc
// @Summary Get voucher by code
// @Description Get voucher details by code
// @Tags voucher
// @Produce json
// @Param code query string true "Voucher code"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /chophimco/api/v1/voucher [get]
func (h *Handler) GetVoucherByCode(ctx *gin.Context) {
	code := ctx.Query("code")
	if code == "" {
		apiwrapper.SendBadRequest(ctx, "Voucher code is required")
		return
	}

	voucher, err := h.voucherUsecase.GetVoucherByCode(ctx, code)
	if err != nil {
		apiwrapper.SendNotFound(ctx, "Voucher not found")
		return
	}

	apiwrapper.SendSuccess(ctx, voucher)
}

// CreateVoucher godoc
// @Summary Create voucher
// @Description Create a new voucher (Admin only)
// @Tags voucher
// @Accept json
// @Produce json
// @Param request body request.CreateVoucher true "Voucher information"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /chophimco/api/v1/voucher/create [post]
func (h *Handler) CreateVoucher(ctx *gin.Context) {
	var req request.CreateVoucher
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	if err := h.voucherUsecase.CreateVoucher(ctx, req); err != nil {
		if err.Error() == "voucher code already exists" {
			apiwrapper.SendBadRequest(ctx, err.Error())
			return
		}
		apiwrapper.SendInternalError(ctx, "Failed to create voucher")
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Voucher created successfully"})
}

// UpdateVoucher godoc
// @Summary Update voucher
// @Description Update voucher information (Admin only)
// @Tags voucher
// @Accept json
// @Produce json
// @Param request body request.UpdateVoucher true "Voucher information"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /chophimco/api/v1/voucher/update [put]
func (h *Handler) UpdateVoucher(ctx *gin.Context) {
	var req request.UpdateVoucher
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	if err := h.voucherUsecase.UpdateVoucher(ctx, req); err != nil {
		apiwrapper.SendInternalError(ctx, "Failed to update voucher")
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Voucher updated successfully"})
}

// DeleteVoucher godoc
// @Summary Delete voucher
// @Description Soft delete a voucher (Admin only)
// @Tags voucher
// @Param id path int true "Voucher ID"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /chophimco/api/v1/voucher/{id} [delete]
func (h *Handler) DeleteVoucher(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid voucher ID")
		return
	}

	if err := h.voucherUsecase.DeleteVoucher(ctx, id); err != nil {
		apiwrapper.SendInternalError(ctx, "Failed to delete voucher")
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Voucher deleted successfully"})
}

// ValidateVoucher godoc
// @Summary Validate voucher
// @Description Validate voucher code
// @Tags voucher
// @Produce json
// @Param code query string true "Voucher code"
// @Param order_value query number true "Order value"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /chophimco/api/v1/voucher/validate [get]
func (h *Handler) ValidateVoucher(ctx *gin.Context) {
	code := ctx.Query("code")
	orderValue, err := strconv.ParseFloat(ctx.Query("order_value"), 64)
	if err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid order value")
		return
	}

	valid, message := h.voucherUsecase.ValidateVoucher(ctx, code, orderValue)

	apiwrapper.SendSuccess(ctx, gin.H{
		"valid":   valid,
		"message": message,
	})
}
