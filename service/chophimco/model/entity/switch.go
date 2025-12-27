package entity

type Switch struct {
	ID    int    `gorm:"primaryKey;column:id;autoIncrement"`
	Name  string `gorm:"column:name;not null"`
	Type  string `gorm:"column:type"` // Linear, Tactile, Clicky
	Brand string `gorm:"column:brand"`
}
