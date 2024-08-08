package models

type Category struct {
	ID   uint   `gorm:"primarykey;autoIncrement" json:"id"`
	Name string `gorm:"column:name;unique;not null" json:"name"`
}
