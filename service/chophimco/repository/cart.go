package repository

import (
	"github.com/leehai1107/chophimco-server/service/chophimco/model/entity"
	"gorm.io/gorm"
)

type ICartRepo interface {
	GetCartByUserID(userID int) (*entity.Cart, error)
	CreateCart(cart *entity.Cart) error
	AddItemToCart(item *entity.CartItem) error
	UpdateCartItem(item *entity.CartItem) error
	RemoveCartItem(itemID int) error
	ClearCart(cartID int) error
}

type cartRepo struct {
	db *gorm.DB
}

func NewCartRepo(db *gorm.DB) ICartRepo {
	return &cartRepo{db: db}
}

func (r *cartRepo) GetCartByUserID(userID int) (*entity.Cart, error) {
	var cart entity.Cart
	err := r.db.Preload("CartItems.ProductVariant.Product").
		Preload("CartItems.ProductVariant.Switch").
		Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepo) CreateCart(cart *entity.Cart) error {
	return r.db.Create(cart).Error
}

func (r *cartRepo) AddItemToCart(item *entity.CartItem) error {
	return r.db.Create(item).Error
}

func (r *cartRepo) UpdateCartItem(item *entity.CartItem) error {
	return r.db.Save(item).Error
}

func (r *cartRepo) RemoveCartItem(itemID int) error {
	return r.db.Delete(&entity.CartItem{}, itemID).Error
}

func (r *cartRepo) ClearCart(cartID int) error {
	return r.db.Where("cart_id = ?", cartID).Delete(&entity.CartItem{}).Error
}
