package http

import (
	"github.com/leehai1107/chophimco-server/service/chophimco/usecase"
)

// IHandler defines all handler interfaces
type IHandler interface {
	IUserHandler
	IProductHandler
	ICartHandler
	IOrderHandler
	IVoucherHandler
	IReviewHandler
}

// Handler implements all handler interfaces
type Handler struct {
	userUsecase    usecase.IUserUsecase
	productUsecase usecase.IProductUsecase
	cartUsecase    usecase.ICartUsecase
	orderUsecase   usecase.IOrderUsecase
	voucherUsecase usecase.IVoucherUsecase
	reviewUsecase  usecase.IReviewUsecase
}

func NewHandler(
	userUsecase usecase.IUserUsecase,
	productUsecase usecase.IProductUsecase,
	cartUsecase usecase.ICartUsecase,
	orderUsecase usecase.IOrderUsecase,
	voucherUsecase usecase.IVoucherUsecase,
	reviewUsecase usecase.IReviewUsecase,
) IHandler {
	return &Handler{
		userUsecase:    userUsecase,
		productUsecase: productUsecase,
		cartUsecase:    cartUsecase,
		orderUsecase:   orderUsecase,
		voucherUsecase: voucherUsecase,
		reviewUsecase:  reviewUsecase,
	}
}
