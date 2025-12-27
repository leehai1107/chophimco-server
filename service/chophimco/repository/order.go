package repository

import (
	"github.com/leehai1107/chophimco-server/service/chophimco/model/entity"
	"gorm.io/gorm"
)

type IOrderRepo interface {
	GetOrderByID(id int) (*entity.Order, error)
	GetOrdersByUserID(userID int) ([]entity.Order, error)
	CreateOrder(order *entity.Order) error
	UpdateOrderStatus(orderID int, status string) error
	CreateOrderItems(items []entity.OrderItem) error
}

type orderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) IOrderRepo {
	return &orderRepo{db: db}
}

func (r *orderRepo) GetOrderByID(id int) (*entity.Order, error) {
	var order entity.Order
	err := r.db.Preload("OrderItems.ProductVariant.Product").
		Preload("OrderItems.ProductVariant.Switch").
		Preload("Voucher").
		Preload("Payment").
		Where("id = ?", id).First(&order).Error
	return &order, err
}

func (r *orderRepo) GetOrdersByUserID(userID int) ([]entity.Order, error) {
	var orders []entity.Order
	err := r.db.Preload("OrderItems.ProductVariant.Product").
		Preload("Voucher").
		Preload("Payment").
		Where("user_id = ?", userID).Order("created_at DESC").Find(&orders).Error
	return orders, err
}

func (r *orderRepo) CreateOrder(order *entity.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepo) UpdateOrderStatus(orderID int, status string) error {
	return r.db.Model(&entity.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

func (r *orderRepo) CreateOrderItems(items []entity.OrderItem) error {
	return r.db.Create(&items).Error
}
