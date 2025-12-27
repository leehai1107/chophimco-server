package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/leehai1107/chophimco-server/pkg/logger"
	"github.com/leehai1107/chophimco-server/service/chophimco/model/entity"
	"github.com/leehai1107/chophimco-server/service/chophimco/model/request"
	"github.com/leehai1107/chophimco-server/service/chophimco/model/response"
	"github.com/leehai1107/chophimco-server/service/chophimco/repository"
)

type IOrderUsecase interface {
	CreateOrder(ctx context.Context, userID int, req request.CreateOrder) (*response.OrderResponse, error)
	GetOrderByID(ctx context.Context, orderID int) (*response.OrderResponse, error)
	GetUserOrders(ctx context.Context, userID int) ([]response.OrderResponse, error)
	UpdateOrderStatus(ctx context.Context, req request.UpdateOrderStatus) error
}

type orderUsecase struct {
	orderRepo   repository.IOrderRepo
	cartRepo    repository.ICartRepo
	voucherRepo repository.IVoucherRepo
	productRepo repository.IProductRepo
	paymentRepo repository.IPaymentRepo
}

func NewOrderUsecase(
	orderRepo repository.IOrderRepo,
	cartRepo repository.ICartRepo,
	voucherRepo repository.IVoucherRepo,
	productRepo repository.IProductRepo,
	paymentRepo repository.IPaymentRepo,
) IOrderUsecase {
	return &orderUsecase{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		voucherRepo: voucherRepo,
		productRepo: productRepo,
		paymentRepo: paymentRepo,
	}
}

func (u *orderUsecase) CreateOrder(ctx context.Context, userID int, req request.CreateOrder) (*response.OrderResponse, error) {
	log := logger.EnhanceWith(ctx)

	// Get user's cart
	cart, err := u.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return nil, errors.New("cart not found")
	}

	if len(cart.CartItems) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Calculate total
	var totalAmount float64
	for _, item := range cart.CartItems {
		if item.ProductVariant != nil {
			totalAmount += item.ProductVariant.Price * float64(item.Quantity)
		}
	}

	// Apply voucher if provided
	var voucherID *int
	var discountAmount float64
	if req.VoucherCode != "" {
		voucher, err := u.voucherRepo.GetVoucherByCode(req.VoucherCode)
		if err != nil {
			return nil, errors.New("invalid voucher code")
		}

		// Validate voucher
		if !voucher.IsActive {
			return nil, errors.New("voucher is not active")
		}

		now := time.Now()
		if voucher.StartAt != nil && voucher.StartAt.After(now) {
			return nil, errors.New("voucher not yet valid")
		}
		if voucher.EndAt != nil && voucher.EndAt.Before(now) {
			return nil, errors.New("voucher has expired")
		}

		if totalAmount < voucher.MinOrderValue {
			return nil, errors.New("order value does not meet voucher minimum")
		}

		// Calculate discount
		if voucher.DiscountType == "percent" {
			discountAmount = totalAmount * (voucher.DiscountValue / 100)
			if voucher.MaxDiscountValue != nil && discountAmount > *voucher.MaxDiscountValue {
				discountAmount = *voucher.MaxDiscountValue
			}
		} else {
			discountAmount = voucher.DiscountValue
		}

		voucherID = &voucher.ID
	}

	finalAmount := totalAmount - discountAmount

	// Create order
	order := &entity.Order{
		UserID:          userID,
		VoucherID:       voucherID,
		DiscountAmount:  discountAmount,
		TotalAmount:     finalAmount,
		Status:          "pending",
		ShippingAddress: req.ShippingAddress,
		CreatedAt:       time.Now(),
	}

	if err := u.orderRepo.CreateOrder(order); err != nil {
		return nil, err
	}

	// Create order items
	orderItems := make([]entity.OrderItem, 0, len(cart.CartItems))
	for _, item := range cart.CartItems {
		if item.ProductVariant != nil {
			orderItems = append(orderItems, entity.OrderItem{
				OrderID:          order.ID,
				ProductVariantID: item.ProductVariantID,
				Price:            item.ProductVariant.Price,
				Quantity:         item.Quantity,
			})

			// Update stock
			if err := u.productRepo.UpdateVariantStock(item.ProductVariantID, -item.Quantity); err != nil {
				log.Errorw("Failed to update stock", "error", err, "variant_id", item.ProductVariantID)
			}
		}
	}

	if err := u.orderRepo.CreateOrderItems(orderItems); err != nil {
		return nil, err
	}

	// Clear cart
	if err := u.cartRepo.ClearCart(cart.ID); err != nil {
		log.Errorw("Failed to clear cart", "error", err)
	}

	// Increment voucher usage
	if voucherID != nil {
		if err := u.voucherRepo.IncrementUsedCount(*voucherID); err != nil {
			log.Errorw("Failed to increment voucher usage", "error", err)
		}
	}

	// Get created order with all relations
	return u.GetOrderByID(ctx, order.ID)
}

func (u *orderUsecase) GetOrderByID(ctx context.Context, orderID int) (*response.OrderResponse, error) {
	order, err := u.orderRepo.GetOrderByID(orderID)
	if err != nil {
		return nil, err
	}
	return u.mapOrderToResponse(order), nil
}

func (u *orderUsecase) GetUserOrders(ctx context.Context, userID int) ([]response.OrderResponse, error) {
	orders, err := u.orderRepo.GetOrdersByUserID(userID)
	if err != nil {
		return nil, err
	}

	result := make([]response.OrderResponse, 0, len(orders))
	for _, order := range orders {
		result = append(result, *u.mapOrderToResponse(&order))
	}
	return result, nil
}

func (u *orderUsecase) UpdateOrderStatus(ctx context.Context, req request.UpdateOrderStatus) error {
	return u.orderRepo.UpdateOrderStatus(req.OrderID, req.Status)
}

func (u *orderUsecase) mapOrderToResponse(order *entity.Order) *response.OrderResponse {
	resp := &response.OrderResponse{
		ID:              order.ID,
		UserID:          order.UserID,
		DiscountAmount:  order.DiscountAmount,
		TotalAmount:     order.TotalAmount,
		Status:          order.Status,
		ShippingAddress: order.ShippingAddress,
		CreatedAt:       order.CreatedAt,
	}

	if order.Voucher != nil {
		code := order.Voucher.Code
		resp.VoucherCode = &code
	}

	items := make([]response.OrderItemResponse, 0, len(order.OrderItems))
	for _, item := range order.OrderItems {
		if item.ProductVariant != nil && item.ProductVariant.Product != nil {
			switchName := ""
			if item.ProductVariant.Switch != nil {
				switchName = item.ProductVariant.Switch.Name
			}

			items = append(items, response.OrderItemResponse{
				ID:          item.ID,
				ProductName: item.ProductVariant.Product.Name,
				Variant: response.ProductVariantResponse{
					ID:             item.ProductVariant.ID,
					ProductID:      item.ProductVariant.ProductID,
					Switch:         &switchName,
					Layout:         item.ProductVariant.Layout,
					ConnectionType: item.ProductVariant.ConnectionType,
					Hotswap:        item.ProductVariant.Hotswap,
					LedType:        item.ProductVariant.LedType,
					Price:          item.ProductVariant.Price,
					Stock:          item.ProductVariant.Stock,
					SKU:            item.ProductVariant.SKU,
				},
				Price:    item.Price,
				Quantity: item.Quantity,
				SubTotal: item.Price * float64(item.Quantity),
			})
		}
	}
	resp.Items = items

	if order.Payment != nil {
		resp.Payment = &response.PaymentResponse{
			ID:            order.Payment.ID,
			PaymentMethod: order.Payment.PaymentMethod,
			PaymentStatus: order.Payment.PaymentStatus,
			PaidAt:        order.Payment.PaidAt,
		}
	}

	return resp
}
