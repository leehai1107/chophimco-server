package repository

import (
	"time"

	"github.com/leehai1107/chophimco-server/service/chophimco/model/entity"
	"gorm.io/gorm"
)

type IVoucherRepo interface {
	GetVoucherByCode(code string) (*entity.Voucher, error)
	GetVoucherByID(id int) (*entity.Voucher, error)
	GetAllVouchers() ([]entity.Voucher, error)
	GetActiveVouchers() ([]entity.Voucher, error)
	CreateVoucher(voucher *entity.Voucher) error
	UpdateVoucher(voucher *entity.Voucher) error
	DeleteVoucher(id int) error
	IncrementUsedCount(voucherID int) error

	// User voucher operations
	GetUserVoucher(userID, voucherID int) (*entity.UserVoucher, error)
	CreateUserVoucher(userVoucher *entity.UserVoucher) error
	IncrementUserVoucherCount(userID, voucherID int) error
}

type voucherRepo struct {
	db *gorm.DB
}

func NewVoucherRepo(db *gorm.DB) IVoucherRepo {
	return &voucherRepo{db: db}
}

func (r *voucherRepo) GetVoucherByCode(code string) (*entity.Voucher, error) {
	var voucher entity.Voucher
	err := r.db.Where("code = ?", code).First(&voucher).Error
	return &voucher, err
}

func (r *voucherRepo) GetVoucherByID(id int) (*entity.Voucher, error) {
	var voucher entity.Voucher
	err := r.db.Where("id = ?", id).First(&voucher).Error
	return &voucher, err
}

func (r *voucherRepo) GetAllVouchers() ([]entity.Voucher, error) {
	var vouchers []entity.Voucher
	err := r.db.Find(&vouchers).Error
	return vouchers, err
}

func (r *voucherRepo) GetActiveVouchers() ([]entity.Voucher, error) {
	var vouchers []entity.Voucher
	now := time.Now()
	err := r.db.Where("is_active = ? AND (start_at IS NULL OR start_at <= ?) AND (end_at IS NULL OR end_at >= ?)",
		true, now, now).Find(&vouchers).Error
	return vouchers, err
}

func (r *voucherRepo) CreateVoucher(voucher *entity.Voucher) error {
	return r.db.Create(voucher).Error
}

func (r *voucherRepo) UpdateVoucher(voucher *entity.Voucher) error {
	return r.db.Save(voucher).Error
}

func (r *voucherRepo) DeleteVoucher(id int) error {
	return r.db.Model(&entity.Voucher{}).Where("id = ?", id).Update("is_active", false).Error
}

func (r *voucherRepo) IncrementUsedCount(voucherID int) error {
	return r.db.Model(&entity.Voucher{}).
		Where("id = ?", voucherID).
		UpdateColumn("used_count", gorm.Expr("used_count + ?", 1)).Error
}

func (r *voucherRepo) GetUserVoucher(userID, voucherID int) (*entity.UserVoucher, error) {
	var userVoucher entity.UserVoucher
	err := r.db.Where("user_id = ? AND voucher_id = ?", userID, voucherID).First(&userVoucher).Error
	return &userVoucher, err
}

func (r *voucherRepo) CreateUserVoucher(userVoucher *entity.UserVoucher) error {
	return r.db.Create(userVoucher).Error
}

func (r *voucherRepo) IncrementUserVoucherCount(userID, voucherID int) error {
	return r.db.Model(&entity.UserVoucher{}).
		Where("user_id = ? AND voucher_id = ?", userID, voucherID).
		UpdateColumn("used_count", gorm.Expr("used_count + ?", 1)).Error
}
