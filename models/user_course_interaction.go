package models

import (
	"gorm.io/gorm"
	"time"
)

type UserCourseInteraction struct {
	gorm.Model
	LastInteraction time.Time `gorm:"column:last_interaction;"`
	TimeSpent       uint      `gorm:"column:time_spent;default:0"`
	Click           bool      `gorm:"column:click;default:false"`
	Enroll          bool      `gorm:"column:enroll;default:false"`
	UserID          uint      `gorm:"column:user_id;not null"`
	CourseID        uint      `gorm:"column:course_id;not null"`

	// Association
	User   User   `gorm:"foreignKey:UserID;references:ID"`
	Course Course `gorm:"foreignKey:CourseID;references:ID"`
}
