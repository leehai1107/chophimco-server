package entity

import (
	"time"
)

type Review struct {
	ID        int       `gorm:"primaryKey;column:id;autoIncrement"`
	UserID    int       `gorm:"column:user_id;not null"`
	ProductID int       `gorm:"column:product_id;not null"`
	Rating    int       `gorm:"column:rating;check:rating >= 1 AND rating <= 5"`
	Comment   string    `gorm:"column:comment;type:text"`
	CreatedAt time.Time `gorm:"column:created_at;default:now()"`

	// Relations
	User    *User    `gorm:"foreignKey:UserID;references:ID"`
	Product *Product `gorm:"foreignKey:ProductID;references:ID"`
}
