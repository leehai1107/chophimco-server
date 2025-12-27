package entity

type Brand struct {
	ID   int    `gorm:"primaryKey;column:id;autoIncrement"`
	Name string `gorm:"column:name;unique;not null"`
}
