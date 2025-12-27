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

type IVoucherUsecase interface {
	GetAllVouchers(ctx context.Context) ([]response.VoucherResponse, error)
	GetActiveVouchers(ctx context.Context) ([]response.VoucherResponse, error)
	GetVoucherByCode(ctx context.Context, code string) (*response.VoucherResponse, error)
	CreateVoucher(ctx context.Context, req request.CreateVoucher) error
	UpdateVoucher(ctx context.Context, req request.UpdateVoucher) error
	DeleteVoucher(ctx context.Context, id int) error
	ValidateVoucher(ctx context.Context, code string, orderValue float64) (bool, string)
}

type voucherUsecase struct {
	repo repository.IVoucherRepo
}

func NewVoucherUsecase(repo repository.IVoucherRepo) IVoucherUsecase {
	return &voucherUsecase{repo: repo}
}

func (u *voucherUsecase) GetAllVouchers(ctx context.Context) ([]response.VoucherResponse, error) {
	vouchers, err := u.repo.GetAllVouchers()
	if err != nil {
		return nil, err
	}
	return u.mapVouchersToResponse(vouchers), nil
}

func (u *voucherUsecase) GetActiveVouchers(ctx context.Context) ([]response.VoucherResponse, error) {
	vouchers, err := u.repo.GetActiveVouchers()
	if err != nil {
		return nil, err
	}
	return u.mapVouchersToResponse(vouchers), nil
}

func (u *voucherUsecase) GetVoucherByCode(ctx context.Context, code string) (*response.VoucherResponse, error) {
	voucher, err := u.repo.GetVoucherByCode(code)
	if err != nil {
		return nil, err
	}
	resp := u.mapVoucherToResponse(voucher)
	return &resp, nil
}

func (u *voucherUsecase) CreateVoucher(ctx context.Context, req request.CreateVoucher) error {
	// Check if code already exists
	_, err := u.repo.GetVoucherByCode(req.Code)
	if err == nil {
		return errors.New("voucher code already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	voucher := &entity.Voucher{
		Code:             req.Code,
		Description:      req.Description,
		DiscountType:     req.DiscountType,
		DiscountValue:    req.DiscountValue,
		MinOrderValue:    req.MinOrderValue,
		MaxDiscountValue: req.MaxDiscountValue,
		UsageLimit:       req.UsageLimit,
		UsagePerUser:     req.UsagePerUser,
		StartAt:          req.StartAt,
		EndAt:            req.EndAt,
		IsActive:         true,
		CreatedAt:        time.Now(),
	}

	return u.repo.CreateVoucher(voucher)
}

func (u *voucherUsecase) UpdateVoucher(ctx context.Context, req request.UpdateVoucher) error {
	voucher, err := u.repo.GetVoucherByID(req.ID)
	if err != nil {
		return err
	}

	if req.Description != "" {
		voucher.Description = req.Description
	}
	if req.DiscountType != "" {
		voucher.DiscountType = req.DiscountType
	}
	if req.DiscountValue > 0 {
		voucher.DiscountValue = req.DiscountValue
	}
	if req.MinOrderValue > 0 {
		voucher.MinOrderValue = req.MinOrderValue
	}
	if req.MaxDiscountValue != nil {
		voucher.MaxDiscountValue = req.MaxDiscountValue
	}
	if req.UsageLimit != nil {
		voucher.UsageLimit = req.UsageLimit
	}
	if req.UsagePerUser != nil {
		voucher.UsagePerUser = *req.UsagePerUser
	}
	if req.StartAt != nil {
		voucher.StartAt = req.StartAt
	}
	if req.EndAt != nil {
		voucher.EndAt = req.EndAt
	}
	if req.IsActive != nil {
		voucher.IsActive = *req.IsActive
	}

	return u.repo.UpdateVoucher(voucher)
}

func (u *voucherUsecase) DeleteVoucher(ctx context.Context, id int) error {
	return u.repo.DeleteVoucher(id)
}

func (u *voucherUsecase) ValidateVoucher(ctx context.Context, code string, orderValue float64) (bool, string) {
	voucher, err := u.repo.GetVoucherByCode(code)
	if err != nil {
		return false, "Invalid voucher code"
	}

	if !voucher.IsActive {
		return false, "Voucher is not active"
	}

	now := time.Now()
	if voucher.StartAt != nil && voucher.StartAt.After(now) {
		return false, "Voucher not yet valid"
	}
	if voucher.EndAt != nil && voucher.EndAt.Before(now) {
		return false, "Voucher has expired"
	}

	if orderValue < voucher.MinOrderValue {
		return false, "Order value does not meet minimum requirement"
	}

	if voucher.UsageLimit != nil && voucher.UsedCount >= *voucher.UsageLimit {
		return false, "Voucher usage limit reached"
	}

	return true, "Voucher is valid"
}

func (u *voucherUsecase) mapVouchersToResponse(vouchers []entity.Voucher) []response.VoucherResponse {
	result := make([]response.VoucherResponse, 0, len(vouchers))
	for _, v := range vouchers {
		result = append(result, u.mapVoucherToResponse(&v))
	}
	return result
}

func (u *voucherUsecase) mapVoucherToResponse(voucher *entity.Voucher) response.VoucherResponse {
	return response.VoucherResponse{
		ID:               voucher.ID,
		Code:             voucher.Code,
		Description:      voucher.Description,
		DiscountType:     voucher.DiscountType,
		DiscountValue:    voucher.DiscountValue,
		MinOrderValue:    voucher.MinOrderValue,
		MaxDiscountValue: voucher.MaxDiscountValue,
		UsageLimit:       voucher.UsageLimit,
		UsagePerUser:     voucher.UsagePerUser,
		UsedCount:        voucher.UsedCount,
		StartAt:          voucher.StartAt,
		EndAt:            voucher.EndAt,
		IsActive:         voucher.IsActive,
		CreatedAt:        voucher.CreatedAt,
	}
}
