package entity

type Role struct {
	ID   int    `gorm:"primaryKey;column:id"`
	Name string `gorm:"column:name;unique;not null"`
}
