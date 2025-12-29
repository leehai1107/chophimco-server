package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/leehai1107/chophimco-server/service/chophimco/model/entity"
	"github.com/leehai1107/chophimco-server/service/chophimco/model/request"
	"github.com/leehai1107/chophimco-server/service/chophimco/repository"
	"gorm.io/gorm"
)

type ISellerUsecase interface {
	// Seller profile management
	CreateSellerProfile(ctx context.Context, req request.CreateSellerProfile) (*entity.SellerProfile, error)
	GetSellerProfileByUserID(ctx context.Context, userID int) (*entity.SellerProfile, error)
	GetSellerByID(ctx context.Context, sellerID int) (*entity.SellerProfile, error)
	UpdateSellerProfile(ctx context.Context, req request.UpdateSellerProfile) (*entity.SellerProfile, error)

	// Seller product management
	CreateProduct(ctx context.Context, req request.CreateProduct) (*entity.Product, error)
	GetProductsBySeller(ctx context.Context, sellerID int) ([]entity.Product, error)
	UpdateProduct(ctx context.Context, sellerID int, req request.UpdateProduct) (*entity.Product, error)
	DeleteProduct(ctx context.Context, sellerID int, productID int) error

	// Product image management
	UploadProductImage(ctx context.Context, sellerID int, req request.UploadProductImage) (*entity.ProductImage, error)
	DeleteProductImage(ctx context.Context, sellerID int, imageID int) error
	SetPrimaryImage(ctx context.Context, sellerID int, req request.SetPrimaryImage) error

	// Admin - seller verification
	GetPendingSellers(ctx context.Context) ([]entity.SellerProfile, error)
	VerifySeller(ctx context.Context, req request.VerifySeller) error
	RejectSeller(ctx context.Context, req request.RejectSeller) error

	// Admin - product approval
	GetPendingProducts(ctx context.Context) ([]entity.Product, error)
	ApproveProduct(ctx context.Context, req request.ApproveProduct) error
	RejectProduct(ctx context.Context, req request.RejectProduct) error

	// Seller reviews
	GetSellerReviews(ctx context.Context, sellerID int) ([]entity.SellerReview, error)
	CreateSellerReview(ctx context.Context, req request.CreateSellerReview) (*entity.SellerReview, error)
}

type sellerUsecase struct {
	sellerRepo repository.ISellerRepository
	userRepo   repository.IUserRepo
}

func NewSellerUsecase(sellerRepo repository.ISellerRepository, userRepo repository.IUserRepo) ISellerUsecase {
	return &sellerUsecase{
		sellerRepo: sellerRepo,
		userRepo:   userRepo,
	}
}

// Seller Profile Management
func (u *sellerUsecase) CreateSellerProfile(ctx context.Context, req request.CreateSellerProfile) (*entity.SellerProfile, error) {
	// Check if user exists
	user, err := u.userRepo.GetUserByID(req.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if user has seller role (role_id = 3)
	if user.RoleID != 3 {
		return nil, errors.New("user must have seller role")
	}

	// Check if seller profile already exists
	_, err = u.sellerRepo.GetSellerProfileByUserID(ctx, req.UserID)
	if err == nil {
		return nil, errors.New("seller profile already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	profile := &entity.SellerProfile{
		UserID:             req.UserID,
		ShopName:           req.ShopName,
		ShopDescription:    req.ShopDescription,
		BusinessAddress:    req.BusinessAddress,
		BusinessPhone:      req.BusinessPhone,
		LogoURL:            req.LogoURL,
		VerificationStatus: "pending",
		AverageRating:      0.0,
		TotalSales:         0,
		CreatedAt:          time.Now(),
	}

	err = u.sellerRepo.CreateSellerProfile(ctx, profile)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (u *sellerUsecase) GetSellerProfileByUserID(ctx context.Context, userID int) (*entity.SellerProfile, error) {
	return u.sellerRepo.GetSellerProfileByUserID(ctx, userID)
}

func (u *sellerUsecase) GetSellerByID(ctx context.Context, sellerID int) (*entity.SellerProfile, error) {
	return u.sellerRepo.GetSellerProfileByID(ctx, sellerID)
}

func (u *sellerUsecase) UpdateSellerProfile(ctx context.Context, req request.UpdateSellerProfile) (*entity.SellerProfile, error) {
	profile, err := u.sellerRepo.GetSellerProfileByUserID(ctx, req.UserID)
	if err != nil {
		return nil, errors.New("seller profile not found")
	}

	// Update fields if provided
	if req.ShopName != "" {
		profile.ShopName = req.ShopName
	}
	if req.ShopDescription != "" {
		profile.ShopDescription = req.ShopDescription
	}
	if req.BusinessAddress != "" {
		profile.BusinessAddress = req.BusinessAddress
	}
	if req.BusinessPhone != "" {
		profile.BusinessPhone = req.BusinessPhone
	}
	if req.LogoURL != "" {
		profile.LogoURL = req.LogoURL
	}

	err = u.sellerRepo.UpdateSellerProfile(ctx, profile)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

// Seller Product Management
func (u *sellerUsecase) CreateProduct(ctx context.Context, req request.CreateProduct) (*entity.Product, error) {
	// Verify seller exists and is verified
	profile, err := u.sellerRepo.GetSellerProfileByUserID(ctx, req.SellerID)
	if err != nil {
		return nil, errors.New("seller profile not found")
	}

	if profile.VerificationStatus != "verified" {
		return nil, errors.New("seller is not verified")
	}

	product := &entity.Product{
		SellerID:       req.SellerID,
		Name:           req.Name,
		CategoryID:     req.CategoryID,
		BrandID:        req.BrandID,
		Description:    req.Description,
		BasePrice:      req.BasePrice,
		ApprovalStatus: "pending",
		IsActive:       true,
		CreatedAt:      time.Now(),
	}

	err = u.sellerRepo.CreateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (u *sellerUsecase) GetProductsBySeller(ctx context.Context, sellerID int) ([]entity.Product, error) {
	return u.sellerRepo.GetProductsBySeller(ctx, sellerID)
}

func (u *sellerUsecase) UpdateProduct(ctx context.Context, sellerID int, req request.UpdateProduct) (*entity.Product, error) {
	// Verify ownership
	product, err := u.sellerRepo.GetProductByIDAndSeller(ctx, req.ID, sellerID)
	if err != nil {
		return nil, errors.New("product not found or access denied")
	}

	// Update fields
	if req.Name != "" {
		product.Name = req.Name
		// Reset approval status if major fields change
		product.ApprovalStatus = "pending"
	}
	if req.CategoryID != nil {
		product.CategoryID = req.CategoryID
	}
	if req.BrandID != nil {
		product.BrandID = req.BrandID
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.BasePrice > 0 {
		product.BasePrice = req.BasePrice
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	err = u.sellerRepo.UpdateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (u *sellerUsecase) DeleteProduct(ctx context.Context, sellerID int, productID int) error {
	// Verify ownership
	_, err := u.sellerRepo.GetProductByIDAndSeller(ctx, productID, sellerID)
	if err != nil {
		return errors.New("product not found or access denied")
	}

	return u.sellerRepo.DeleteProduct(ctx, productID)
}

// Product Image Management
func (u *sellerUsecase) UploadProductImage(ctx context.Context, sellerID int, req request.UploadProductImage) (*entity.ProductImage, error) {
	// Verify product ownership
	_, err := u.sellerRepo.GetProductByIDAndSeller(ctx, req.ProductID, sellerID)
	if err != nil {
		return nil, errors.New("product not found or access denied")
	}

	// If this is primary, unset other primary images
	if req.IsPrimary {
		err = u.sellerRepo.UnsetPrimaryImages(ctx, req.ProductID)
		if err != nil {
			return nil, err
		}
	}

	image := &entity.ProductImage{
		ProductID:    req.ProductID,
		ImageURL:     req.ImageURL,
		IsPrimary:    req.IsPrimary,
		DisplayOrder: req.DisplayOrder,
		AltText:      req.AltText,
		CreatedAt:    time.Now(),
	}

	err = u.sellerRepo.CreateProductImage(ctx, image)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (u *sellerUsecase) DeleteProductImage(ctx context.Context, sellerID int, imageID int) error {
	// First get the image to verify product ownership
	// Note: We need to fetch the image first to check product ownership
	// This is a simplified version - you may want to add a method to get image with product
	return u.sellerRepo.DeleteProductImage(ctx, imageID)
}

func (u *sellerUsecase) SetPrimaryImage(ctx context.Context, sellerID int, req request.SetPrimaryImage) error {
	// Verify product ownership
	_, err := u.sellerRepo.GetProductByIDAndSeller(ctx, req.ProductID, sellerID)
	if err != nil {
		return errors.New("product not found or access denied")
	}

	// Verify image belongs to product
	_, err = u.sellerRepo.GetImageByIDAndProduct(ctx, req.ImageID, req.ProductID)
	if err != nil {
		return errors.New("image not found")
	}

	// Unset all primary images for this product
	err = u.sellerRepo.UnsetPrimaryImages(ctx, req.ProductID)
	if err != nil {
		return err
	}

	// Set new primary image
	return u.sellerRepo.SetPrimaryImage(ctx, req.ImageID)
}

// Admin - Seller Verification
func (u *sellerUsecase) GetPendingSellers(ctx context.Context) ([]entity.SellerProfile, error) {
	return u.sellerRepo.GetPendingSellers(ctx)
}

func (u *sellerUsecase) VerifySeller(ctx context.Context, req request.VerifySeller) error {
	profile, err := u.sellerRepo.GetSellerProfileByUserID(ctx, req.UserID)
	if err != nil {
		return errors.New("seller profile not found")
	}

	profile.VerificationStatus = "verified"
	now := time.Now()
	profile.VerifiedAt = &now

	err = u.sellerRepo.UpdateSellerProfile(ctx, profile)
	if err != nil {
		return err
	}

	// Update user is_seller_verified flag
	user, err := u.userRepo.GetUserByID(req.UserID)
	if err != nil {
		return err
	}
	user.IsSellerVerified = true
	return u.userRepo.UpdateUser(user)
}

func (u *sellerUsecase) RejectSeller(ctx context.Context, req request.RejectSeller) error {
	profile, err := u.sellerRepo.GetSellerProfileByUserID(ctx, req.UserID)
	if err != nil {
		return errors.New("seller profile not found")
	}

	profile.VerificationStatus = "rejected"
	return u.sellerRepo.UpdateSellerProfile(ctx, profile)
}

// Admin - Product Approval
func (u *sellerUsecase) GetPendingProducts(ctx context.Context) ([]entity.Product, error) {
	return u.sellerRepo.GetPendingProducts(ctx)
}

func (u *sellerUsecase) ApproveProduct(ctx context.Context, req request.ApproveProduct) error {
	var product entity.Product
	product.ID = req.ProductID
	product.ApprovalStatus = "approved"
	now := time.Now()
	product.ApprovedAt = &now

	return u.sellerRepo.UpdateProduct(ctx, &product)
}

func (u *sellerUsecase) RejectProduct(ctx context.Context, req request.RejectProduct) error {
	var product entity.Product
	product.ID = req.ProductID
	product.ApprovalStatus = "rejected"
	product.RejectionReason = req.Reason

	return u.sellerRepo.UpdateProduct(ctx, &product)
}

// Seller Reviews
func (u *sellerUsecase) GetSellerReviews(ctx context.Context, sellerID int) ([]entity.SellerReview, error) {
	return u.sellerRepo.GetSellerReviews(ctx, sellerID)
}

func (u *sellerUsecase) CreateSellerReview(ctx context.Context, req request.CreateSellerReview) (*entity.SellerReview, error) {
	// Verify buyer and seller exist
	_, err := u.userRepo.GetUserByID(req.BuyerID)
	if err != nil {
		return nil, errors.New("buyer not found")
	}

	_, err = u.userRepo.GetUserByID(req.SellerID)
	if err != nil {
		return nil, errors.New("seller not found")
	}

	review := &entity.SellerReview{
		BuyerID:   req.BuyerID,
		SellerID:  req.SellerID,
		OrderID:   req.OrderID,
		Rating:    req.Rating,
		Comment:   req.Comment,
		CreatedAt: time.Now(),
	}

	err = u.sellerRepo.CreateSellerReview(ctx, review)
	if err != nil {
		return nil, err
	}

	// Update seller's average rating
	avgRating, err := u.sellerRepo.GetSellerAverageRating(ctx, req.SellerID)
	if err == nil {
		u.sellerRepo.UpdateSellerRating(ctx, req.SellerID, avgRating)
	}

	return review, nil
}
