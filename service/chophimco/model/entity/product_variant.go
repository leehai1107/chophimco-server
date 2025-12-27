package entity

type ProductVariant struct {
	ID             int     `gorm:"primaryKey;column:id;autoIncrement"`
	ProductID      int     `gorm:"column:product_id;not null"`
	SwitchID       *int    `gorm:"column:switch_id"`
	Layout         string  `gorm:"column:layout"`          // 60%, 65%, TKL, Fullsize
	ConnectionType string  `gorm:"column:connection_type"` // Wired, Wireless, Bluetooth
	Hotswap        bool    `gorm:"column:hotswap;default:false"`
	LedType        string  `gorm:"column:led_type"` // RGB, White
	Price          float64 `gorm:"column:price;not null"`
	Stock          int     `gorm:"column:stock;default:0"`
	SKU            string  `gorm:"column:sku;unique"`

	// Relations
	Product *Product `gorm:"foreignKey:ProductID;references:ID"`
	Switch  *Switch  `gorm:"foreignKey:SwitchID;references:ID"`
}
