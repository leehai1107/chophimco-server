package entity

import (
	"time"
)

type SellerReview struct {
	ID        int       `gorm:"primaryKey;column:id;autoIncrement"`
	BuyerID   int       `gorm:"column:buyer_id;not null"`
	SellerID  int       `gorm:"column:seller_id;not null"`
	OrderID   int       `gorm:"column:order_id;not null"`
	Rating    int       `gorm:"column:rating;check:rating >= 1 AND rating <= 5"`
	Comment   string    `gorm:"column:comment;type:text"`
	CreatedAt time.Time `gorm:"column:created_at;default:now()"`

	// Relations
	Buyer  *User  `gorm:"foreignKey:BuyerID;references:ID"`
	Seller *User  `gorm:"foreignKey:SellerID;references:ID"`
	Order  *Order `gorm:"foreignKey:OrderID;references:ID"`
}
