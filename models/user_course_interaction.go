package models

import (
	"gorm.io/gorm"
	"time"
)

type UserCourseInteraction struct {
	ID              uint `gorm:"primarykey;auto_increment" json:"id"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
	LastInteraction time.Time      `gorm:"column:last_interaction;" json:"last_interaction"`
	TimeSpent       uint           `gorm:"column:time_spent;default:0" json:"time_spent"`
	Click           bool           `gorm:"column:click;default:false" json:"click"`
	Enroll          bool           `gorm:"column:enroll;default:false" json:"enroll"`
	UserID          uint           `gorm:"column:user_id;not null" json:"user_id"`
	CourseID        uint           `gorm:"column:course_id;not null" json:"course_id"`

	// Association
	User   User   `gorm:"foreignKey:UserID;references:ID"`
	Course Course `gorm:"foreignKey:CourseID;references:ID"`
}
