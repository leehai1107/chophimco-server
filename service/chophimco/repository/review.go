package repository

import (
	"github.com/leehai1107/chophimco-server/service/chophimco/model/entity"
	"gorm.io/gorm"
)

type IReviewRepo interface {
	GetReviewsByProductID(productID int) ([]entity.Review, error)
	GetReviewByUserAndProduct(userID, productID int) (*entity.Review, error)
	CreateReview(review *entity.Review) error
	UpdateReview(review *entity.Review) error
	DeleteReview(id int) error
}

type reviewRepo struct {
	db *gorm.DB
}

func NewReviewRepo(db *gorm.DB) IReviewRepo {
	return &reviewRepo{db: db}
}

func (r *reviewRepo) GetReviewsByProductID(productID int) ([]entity.Review, error) {
	var reviews []entity.Review
	err := r.db.Preload("User").Where("product_id = ?", productID).
		Order("created_at DESC").Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepo) GetReviewByUserAndProduct(userID, productID int) (*entity.Review, error) {
	var review entity.Review
	err := r.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&review).Error
	return &review, err
}

func (r *reviewRepo) CreateReview(review *entity.Review) error {
	return r.db.Create(review).Error
}

func (r *reviewRepo) UpdateReview(review *entity.Review) error {
	return r.db.Save(review).Error
}

func (r *reviewRepo) DeleteReview(id int) error {
	return r.db.Delete(&entity.Review{}, id).Error
}
