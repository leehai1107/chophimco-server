package repository

import (
	"context"

	"github.com/leehai1107/chophimco-server/service/chophimco/model/entity"
	"gorm.io/gorm"
)

type ISellerRepository interface {
	// Seller profile
	CreateSellerProfile(ctx context.Context, profile *entity.SellerProfile) error
	GetSellerProfileByUserID(ctx context.Context, userID int) (*entity.SellerProfile, error)
	GetSellerProfileByID(ctx context.Context, id int) (*entity.SellerProfile, error)
	UpdateSellerProfile(ctx context.Context, profile *entity.SellerProfile) error
	GetPendingSellers(ctx context.Context) ([]entity.SellerProfile, error)

	// Product management
	CreateProduct(ctx context.Context, product *entity.Product) error
	GetProductsBySeller(ctx context.Context, sellerID int) ([]entity.Product, error)
	GetProductByIDAndSeller(ctx context.Context, productID int, sellerID int) (*entity.Product, error)
	UpdateProduct(ctx context.Context, product *entity.Product) error
	DeleteProduct(ctx context.Context, productID int) error
	GetPendingProducts(ctx context.Context) ([]entity.Product, error)

	// Product images
	CreateProductImage(ctx context.Context, image *entity.ProductImage) error
	GetImageByIDAndProduct(ctx context.Context, imageID int, productID int) (*entity.ProductImage, error)
	DeleteProductImage(ctx context.Context, imageID int) error
	UnsetPrimaryImages(ctx context.Context, productID int) error
	SetPrimaryImage(ctx context.Context, imageID int) error

	// Seller reviews
	GetSellerReviews(ctx context.Context, sellerID int) ([]entity.SellerReview, error)
	CreateSellerReview(ctx context.Context, review *entity.SellerReview) error
	GetSellerAverageRating(ctx context.Context, sellerID int) (float64, error)
	UpdateSellerRating(ctx context.Context, sellerID int, rating float64) error
}

type sellerRepo struct {
	db *gorm.DB
}

func NewSellerRepo(db *gorm.DB) ISellerRepository {
	return &sellerRepo{db: db}
}

// Seller Profile Methods
func (r *sellerRepo) CreateSellerProfile(ctx context.Context, profile *entity.SellerProfile) error {
	return r.db.WithContext(ctx).Create(profile).Error
}

func (r *sellerRepo) GetSellerProfileByUserID(ctx context.Context, userID int) (*entity.SellerProfile, error) {
	var profile entity.SellerProfile
	err := r.db.WithContext(ctx).Preload("User").Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *sellerRepo) GetSellerProfileByID(ctx context.Context, id int) (*entity.SellerProfile, error) {
	var profile entity.SellerProfile
	err := r.db.WithContext(ctx).Preload("User").Where("id = ?", id).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *sellerRepo) UpdateSellerProfile(ctx context.Context, profile *entity.SellerProfile) error {
	return r.db.WithContext(ctx).Save(profile).Error
}

func (r *sellerRepo) GetPendingSellers(ctx context.Context) ([]entity.SellerProfile, error) {
	var sellers []entity.SellerProfile
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("verification_status = ?", "pending").
		Order("created_at DESC").
		Find(&sellers).Error
	return sellers, err
}

// Product Management Methods
func (r *sellerRepo) CreateProduct(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *sellerRepo) GetProductsBySeller(ctx context.Context, sellerID int) ([]entity.Product, error) {
	var products []entity.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Brand").
		Preload("Variants").
		Preload("Images").
		Where("seller_id = ?", sellerID).
		Order("created_at DESC").
		Find(&products).Error
	return products, err
}

func (r *sellerRepo) GetProductByIDAndSeller(ctx context.Context, productID int, sellerID int) (*entity.Product, error) {
	var product entity.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Brand").
		Preload("Variants").
		Preload("Images").
		Where("id = ? AND seller_id = ?", productID, sellerID).
		First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *sellerRepo) UpdateProduct(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *sellerRepo) DeleteProduct(ctx context.Context, productID int) error {
	return r.db.WithContext(ctx).Delete(&entity.Product{}, productID).Error
}

func (r *sellerRepo) GetPendingProducts(ctx context.Context) ([]entity.Product, error) {
	var products []entity.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Brand").
		Preload("Seller").
		Preload("Images").
		Where("approval_status = ?", "pending").
		Order("created_at DESC").
		Find(&products).Error
	return products, err
}

// Product Image Methods
func (r *sellerRepo) CreateProductImage(ctx context.Context, image *entity.ProductImage) error {
	return r.db.WithContext(ctx).Create(image).Error
}

func (r *sellerRepo) GetImageByIDAndProduct(ctx context.Context, imageID int, productID int) (*entity.ProductImage, error) {
	var image entity.ProductImage
	err := r.db.WithContext(ctx).
		Where("id = ? AND product_id = ?", imageID, productID).
		First(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

func (r *sellerRepo) DeleteProductImage(ctx context.Context, imageID int) error {
	return r.db.WithContext(ctx).Delete(&entity.ProductImage{}, imageID).Error
}

func (r *sellerRepo) UnsetPrimaryImages(ctx context.Context, productID int) error {
	return r.db.WithContext(ctx).
		Model(&entity.ProductImage{}).
		Where("product_id = ?", productID).
		Update("is_primary", false).Error
}

func (r *sellerRepo) SetPrimaryImage(ctx context.Context, imageID int) error {
	return r.db.WithContext(ctx).
		Model(&entity.ProductImage{}).
		Where("id = ?", imageID).
		Update("is_primary", true).Error
}

// Seller Review Methods
func (r *sellerRepo) GetSellerReviews(ctx context.Context, sellerID int) ([]entity.SellerReview, error) {
	var reviews []entity.SellerReview
	err := r.db.WithContext(ctx).
		Preload("Buyer").
		Preload("Order").
		Where("seller_id = ?", sellerID).
		Order("created_at DESC").
		Find(&reviews).Error
	return reviews, err
}

func (r *sellerRepo) CreateSellerReview(ctx context.Context, review *entity.SellerReview) error {
	return r.db.WithContext(ctx).Create(review).Error
}

func (r *sellerRepo) GetSellerAverageRating(ctx context.Context, sellerID int) (float64, error) {
	var avgRating float64
	err := r.db.WithContext(ctx).
		Model(&entity.SellerReview{}).
		Where("seller_id = ?", sellerID).
		Select("COALESCE(AVG(rating), 0)").
		Scan(&avgRating).Error
	return avgRating, err
}

func (r *sellerRepo) UpdateSellerRating(ctx context.Context, sellerID int, rating float64) error {
	return r.db.WithContext(ctx).
		Model(&entity.SellerProfile{}).
		Where("user_id = ?", sellerID).
		Update("average_rating", rating).Error
}
