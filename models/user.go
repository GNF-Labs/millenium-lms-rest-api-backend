package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name  string  `gorm:"column:name;size:64;unique;not null"`
	Email string  `gorm:"column:email;size:64;unique;not null"`
	Age   uint8   `gorm:"column:age;not null"`
	About *string `gorm:"column:about;type:text"`
	Role  string  `gorm:"column:role;type:varchar(32);not null;check:role IN ('admin', 'student')'"`
}
