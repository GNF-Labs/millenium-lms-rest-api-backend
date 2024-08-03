package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name string `gorm:"column:name;unique;not null"`
}
