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
	Viewed          bool           `gorm:"column:viewed;default:false" json:"viewed"`
	Registered      bool           `gorm:"column:registered;default:false" json:"registered"`
	Completed       bool           `gorm:"column:completed;default:false" json:"completed"`
	CompletionRate  float32        `gorm:"column:completion_rate;check:completion_rate >= 0 AND completion_rate <= 1" json:"completion_rate"`
	UserID          uint           `gorm:"column:user_id;not null" json:"user_id"`
	CourseID        uint           `gorm:"column:course_id;not null" json:"course_id"`

	// Association
	User   User   `gorm:"foreignKey:UserID;references:ID"`
	Course Course `gorm:"foreignKey:CourseID;references:ID"`
}
