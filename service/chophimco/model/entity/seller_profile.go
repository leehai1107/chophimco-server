package entity

import (
	"time"
)

type SellerProfile struct {
	ID                 int        `gorm:"primaryKey;column:id;autoIncrement"`
	UserID             int        `gorm:"column:user_id;unique;not null"`
	ShopName           string     `gorm:"column:shop_name;not null"`
	ShopDescription    string     `gorm:"column:shop_description;type:text"`
	BusinessAddress    string     `gorm:"column:business_address;type:text"`
	BusinessPhone      string     `gorm:"column:business_phone"`
	LogoURL            string     `gorm:"column:logo_url;type:text"`
	VerificationStatus string     `gorm:"column:verification_status;default:pending"`
	AverageRating      float64    `gorm:"column:average_rating;type:decimal(3,2);default:0.0"`
	TotalSales         int        `gorm:"column:total_sales;default:0"`
	CreatedAt          time.Time  `gorm:"column:created_at;default:now()"`
	VerifiedAt         *time.Time `gorm:"column:verified_at"`

	// Relations
	User *User `gorm:"foreignKey:UserID;references:ID"`
}
