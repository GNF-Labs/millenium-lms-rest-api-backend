package models

import (
	"gorm.io/gorm"
	"time"
)

type Subchapter struct {
	ID        uint `gorm:"primarykey;auto_increment" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"column:name;size:64;not null" json:"name"`
	Content   string         `gorm:"column:content;type:text;not null" json:"content"`
	ChapterID uint           `gorm:"primarykey;column:chapter_id;not null" json:"chapter_id"`
	CourseID  uint           `gorm:"primarykey;column:course_id;not null" json:"course_id"`
	Chapter   Chapter        `gorm:"foreignKey:ChapterID;references:ID"`
	Course    Course         `gorm:"foreignKey:CourseID;references:ID"`
	Order     uint           `gorm:"column:order;default:0" json:"order"`
}
