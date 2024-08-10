package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint   `gorm:"primarykey;auto_increment" json:"id"`
	Username  string `gorm:"size:24;unique;not null;default:''" json:"username"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	FullName  string         `gorm:"column:full_name;size:64;unique;not null;default:''" json:"full_name"`
	Email     string         `gorm:"column:email;size:64;unique;not null" json:"email"`
	About     string         `gorm:"column:about;type:text" json:"about"`
	Role      string         `gorm:"column:role;type:varchar(32);not null;check:role IN ('admin', 'student')" json:"role"`
	Password  string         `gorm:"column:password;size:255;not null;default:''" json:"password"`
	ImageURL  string         `gorm:"column:image_url;" json:"image_url"`
}
