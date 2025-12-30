package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/leehai1107/chophimco-server/pkg/logger"
	"github.com/leehai1107/chophimco-server/pkg/middleware/auth"
	"github.com/leehai1107/chophimco-server/service/chophimco/model/entity"
	"github.com/leehai1107/chophimco-server/service/chophimco/model/request"
	"github.com/leehai1107/chophimco-server/service/chophimco/model/response"
	"github.com/leehai1107/chophimco-server/service/chophimco/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IUserUsecase interface {
	Login(ctx context.Context, req request.Login) (*response.LoginResponse, error)
	Register(ctx context.Context, req request.Register) error
	GetUserProfile(ctx context.Context, userID int) (*response.UserResponse, error)
}

type userUsecase struct {
	repo       repository.IUserRepo
	jwtService auth.IJWTService
}

func NewUserUsecase(repo repository.IUserRepo, jwtService auth.IJWTService) IUserUsecase {
	return &userUsecase{repo: repo, jwtService: jwtService}
}

func (u *userUsecase) Login(ctx context.Context, req request.Login) (*response.LoginResponse, error) {
	user, err := u.repo.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Compare password hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Get role name
	roleName := ""
	if user.Role != nil {
		roleName = user.Role.Name
	}

	// Generate JWT token
	token, err := u.jwtService.GenerateToken(user.ID, user.Email, roleName)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &response.LoginResponse{
		Token: token,
		User: response.UserResponse{
			ID:        user.ID,
			FullName:  user.FullName,
			Email:     user.Email,
			Phone:     user.Phone,
			Role:      roleName,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

func (u *userUsecase) Register(ctx context.Context, req request.Register) error {
	logger.EnhanceWith(ctx).Info("Register usecase called")

	// Check if email already exists
	_, err := u.repo.GetUserByEmail(req.Email)
	if err == nil {
		return errors.New("email already registered")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create user entity (default role_id = 2 for customer)
	user := &entity.User{
		RoleID:       2, // customer
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Phone:        req.Phone,
		CreatedAt:    time.Now(),
	}

	return u.repo.CreateUser(user)
}

func (u *userUsecase) GetUserProfile(ctx context.Context, userID int) (*response.UserResponse, error) {
	user, err := u.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	roleName := ""
	if user.Role != nil {
		roleName = user.Role.Name
	}

	return &response.UserResponse{
		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      roleName,
		CreatedAt: user.CreatedAt,
	}, nil
}
