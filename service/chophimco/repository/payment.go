package repository

import (
	"github.com/leehai1107/chophimco-server/service/chophimco/model/entity"
	"gorm.io/gorm"
)

type IPaymentRepo interface {
	CreatePayment(payment *entity.Payment) error
	GetPaymentByOrderID(orderID int) (*entity.Payment, error)
	UpdatePaymentStatus(paymentID int, status string) error
}

type paymentRepo struct {
	db *gorm.DB
}

func NewPaymentRepo(db *gorm.DB) IPaymentRepo {
	return &paymentRepo{db: db}
}

func (r *paymentRepo) CreatePayment(payment *entity.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepo) GetPaymentByOrderID(orderID int) (*entity.Payment, error) {
	var payment entity.Payment
	err := r.db.Where("order_id = ?", orderID).First(&payment).Error
	return &payment, err
}

func (r *paymentRepo) UpdatePaymentStatus(paymentID int, status string) error {
	return r.db.Model(&entity.Payment{}).Where("id = ?", paymentID).Update("payment_status", status).Error
}
