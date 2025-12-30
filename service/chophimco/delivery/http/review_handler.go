package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/leehai1107/chophimco-server/pkg/apiwrapper"
	"github.com/leehai1107/chophimco-server/pkg/middleware/auth"
	"github.com/leehai1107/chophimco-server/service/chophimco/model/request"
)

type IReviewHandler interface {
	GetProductReviews(ctx *gin.Context)
	CreateReview(ctx *gin.Context)
}

// GetProductReviews godoc
// @Summary Get product reviews
// @Description Get all reviews for a product
// @Tags review
// @Produce json
// @Param product_id query int true "Product ID"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/review [get]
func (h *Handler) GetProductReviews(ctx *gin.Context) {
	productID, err := strconv.Atoi(ctx.Query("product_id"))
	if err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid product ID")
		return
	}

	reviews, err := h.reviewUsecase.GetProductReviews(ctx, productID)
	if err != nil {
		apiwrapper.SendInternalError(ctx, "Failed to get reviews")
		return
	}

	apiwrapper.SendSuccess(ctx, reviews)
}

// CreateReview godoc
// @Summary Create review
// @Description Create a product review
// @Tags review
// @Accept json
// @Produce json
// @Param request body request.CreateReview true "Review information"
// @Success 200 {object} apiwrapper.APIResponse
// @Router /api/v1/review/create [post]
func (h *Handler) CreateReview(ctx *gin.Context) {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		apiwrapper.SendUnauthorized(ctx, "Unauthorized")
		return
	}

	var req request.CreateReview
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	if err := h.reviewUsecase.CreateReview(ctx, userID, req); err != nil {
		if err.Error() == "you have already reviewed this product" {
			apiwrapper.SendBadRequest(ctx, err.Error())
			return
		}
		apiwrapper.SendInternalError(ctx, "Failed to create review")
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Review created successfully"})
}
