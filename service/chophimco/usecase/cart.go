package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/leehai1107/chophimco-server/service/chophimco/model/entity"
	"github.com/leehai1107/chophimco-server/service/chophimco/model/request"
	"github.com/leehai1107/chophimco-server/service/chophimco/model/response"
	"github.com/leehai1107/chophimco-server/service/chophimco/repository"
	"gorm.io/gorm"
)

type ICartUsecase interface {
	GetCart(ctx context.Context, userID int) (*response.CartResponse, error)
	AddToCart(ctx context.Context, userID int, req request.AddToCart) error
	UpdateCartItem(ctx context.Context, req request.UpdateCartItem) error
	RemoveFromCart(ctx context.Context, cartItemID int) error
	ClearCart(ctx context.Context, userID int) error
}

type cartUsecase struct {
	cartRepo    repository.ICartRepo
	productRepo repository.IProductRepo
}

func NewCartUsecase(cartRepo repository.ICartRepo, productRepo repository.IProductRepo) ICartUsecase {
	return &cartUsecase{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (u *cartUsecase) GetCart(ctx context.Context, userID int) (*response.CartResponse, error) {
	cart, err := u.cartRepo.GetCartByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new cart if not exists
			newCart := &entity.Cart{
				UserID:    userID,
				CreatedAt: time.Now(),
			}
			if err := u.cartRepo.CreateCart(newCart); err != nil {
				return nil, err
			}
			return &response.CartResponse{
				ID:        newCart.ID,
				UserID:    userID,
				Items:     []response.CartItemResponse{},
				TotalItem: 0,
				SubTotal:  0,
			}, nil
		}
		return nil, err
	}

	return u.mapCartToResponse(cart), nil
}

func (u *cartUsecase) AddToCart(ctx context.Context, userID int, req request.AddToCart) error {
	// Get or create cart
	cart, err := u.cartRepo.GetCartByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cart = &entity.Cart{
				UserID:    userID,
				CreatedAt: time.Now(),
			}
			if err := u.cartRepo.CreateCart(cart); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// Verify product variant exists and has stock
	variant, err := u.productRepo.GetVariantByID(req.ProductVariantID)
	if err != nil {
		return errors.New("product variant not found")
	}

	if variant.Stock < req.Quantity {
		return errors.New("insufficient stock")
	}

	// Add item to cart
	item := &entity.CartItem{
		CartID:           cart.ID,
		ProductVariantID: req.ProductVariantID,
		Quantity:         req.Quantity,
	}

	return u.cartRepo.AddItemToCart(item)
}

func (u *cartUsecase) UpdateCartItem(ctx context.Context, req request.UpdateCartItem) error {
	// Update cart item quantity
	item := &entity.CartItem{
		ID:       req.CartItemID,
		Quantity: req.Quantity,
	}
	return u.cartRepo.UpdateCartItem(item)
}

func (u *cartUsecase) RemoveFromCart(ctx context.Context, cartItemID int) error {
	return u.cartRepo.RemoveCartItem(cartItemID)
}

func (u *cartUsecase) ClearCart(ctx context.Context, userID int) error {
	cart, err := u.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return err
	}
	return u.cartRepo.ClearCart(cart.ID)
}

func (u *cartUsecase) mapCartToResponse(cart *entity.Cart) *response.CartResponse {
	items := make([]response.CartItemResponse, 0, len(cart.CartItems))
	var subTotal float64

	for _, item := range cart.CartItems {
		if item.ProductVariant != nil && item.ProductVariant.Product != nil {
			itemSubTotal := item.ProductVariant.Price * float64(item.Quantity)
			subTotal += itemSubTotal

			switchName := ""
			if item.ProductVariant.Switch != nil {
				switchName = item.ProductVariant.Switch.Name
			}

			items = append(items, response.CartItemResponse{
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
				Quantity: item.Quantity,
				Price:    item.ProductVariant.Price,
				SubTotal: itemSubTotal,
			})
		}
	}

	return &response.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		Items:     items,
		TotalItem: len(items),
		SubTotal:  subTotal,
	}
}
