package models

import (
	"gorm.io/gorm"
	"time"
)

type Review struct {
	ID        uint `gorm:"primarykey;auto_increment" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Content   string         `gorm:"column:content;not null" json:"content"`
	Timestamp time.Time      `gorm:"column:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"timestamp"`
	UserID    uint           `gorm:"column:user_id;not null" json:"user_id"`
	CourseID  uint           `gorm:"column:course_id;not null" json:"course_id"`
	User      User           `gorm:"foreignKey:UserID;references:ID"`
	Course    Course         `gorm:"foreignKey:CourseID;references:ID"`
}
