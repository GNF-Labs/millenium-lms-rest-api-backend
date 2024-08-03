package models

import (
	"gorm.io/gorm"
	"time"
)

type Review struct {
	gorm.Model
	Content   string    `gorm:"column:content;not null"`
	Timestamp time.Time `gorm:"column:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UserID    uint      `gorm:"column:user_id;not null"`
	CourseID  uint      `gorm:"column:course_id;not null"`
	User      User      `gorm:"foreignKey:UserID;references:ID"`
	Course    Course    `gorm:"foreignKey:CourseID;references:ID"`
}
