package entity

import (
	"time"
)

type User struct {
	ID               int       `gorm:"primaryKey;column:id;autoIncrement"`
	RoleID           int       `gorm:"column:role_id;not null"`
	Email            string    `gorm:"column:email;unique;not null"`
	PasswordHash     string    `gorm:"column:password_hash;not null"`
	FullName         string    `gorm:"column:full_name"`
	Phone            string    `gorm:"column:phone"`
	IsSellerVerified bool      `gorm:"column:is_seller_verified;default:false"`
	CreatedAt        time.Time `gorm:"column:created_at;default:now()"`

	// Relations
	Role          *Role          `gorm:"foreignKey:RoleID;references:ID"`
	SellerProfile *SellerProfile `gorm:"foreignKey:UserID;references:ID"`
}
