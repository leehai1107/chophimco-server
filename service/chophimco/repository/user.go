package repository

import (
	"errors"

	"github.com/leehai1107/chophimco-server/pkg/logger"
	"github.com/leehai1107/chophimco-server/service/chophimco/model/entity"
	"gorm.io/gorm"
)

type IUserRepo interface {
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByID(id int) (*entity.User, error)
	CreateUser(user *entity.User) error
	UpdateUser(user *entity.User) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) IUserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) GetUserByEmail(email string) (*entity.User, error) {
	logger.Info("GetUserByEmail repository method called")
	var user entity.User
	result := r.db.Preload("Role").Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepo) GetUserByID(id int) (*entity.User, error) {
	var user entity.User
	result := r.db.Preload("Role").Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepo) CreateUser(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepo) UpdateUser(user *entity.User) error {
	return r.db.Save(user).Error
}
