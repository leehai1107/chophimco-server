package entity

type Category struct {
	ID   int    `gorm:"primaryKey;column:id;autoIncrement"`
	Name string `gorm:"column:name;not null"`
}
