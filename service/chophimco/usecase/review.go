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

type IReviewUsecase interface {
	GetProductReviews(ctx context.Context, productID int) ([]response.ReviewResponse, error)
	CreateReview(ctx context.Context, userID int, req request.CreateReview) error
}

type reviewUsecase struct {
	reviewRepo repository.IReviewRepo
}

func NewReviewUsecase(reviewRepo repository.IReviewRepo) IReviewUsecase {
	return &reviewUsecase{reviewRepo: reviewRepo}
}

func (u *reviewUsecase) GetProductReviews(ctx context.Context, productID int) ([]response.ReviewResponse, error) {
	reviews, err := u.reviewRepo.GetReviewsByProductID(productID)
	if err != nil {
		return nil, err
	}

	result := make([]response.ReviewResponse, 0, len(reviews))
	for _, review := range reviews {
		resp := response.ReviewResponse{
			ID:        review.ID,
			ProductID: review.ProductID,
			Rating:    review.Rating,
			Comment:   review.Comment,
			CreatedAt: review.CreatedAt,
		}

		if review.User != nil {
			resp.UserName = review.User.FullName
		}

		if review.Product != nil {
			resp.ProductName = review.Product.Name
		}

		result = append(result, resp)
	}

	return result, nil
}

func (u *reviewUsecase) CreateReview(ctx context.Context, userID int, req request.CreateReview) error {
	// Check if user already reviewed this product
	_, err := u.reviewRepo.GetReviewByUserAndProduct(userID, req.ProductID)
	if err == nil {
		return errors.New("you have already reviewed this product")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	review := &entity.Review{
		UserID:    userID,
		ProductID: req.ProductID,
		Rating:    req.Rating,
		Comment:   req.Comment,
		CreatedAt: time.Now(),
	}

	return u.reviewRepo.CreateReview(review)
}
