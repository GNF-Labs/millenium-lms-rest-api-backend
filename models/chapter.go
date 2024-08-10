package models

import (
	"gorm.io/gorm"
	"time"
)

type Chapter struct {
	ID          uint `gorm:"primarykey;auto_increment" json:"id"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string         `gorm:"column:name;type:varchar(50);not null" json:"name"`
	Description string         `gorm:"column:description;type:text" json:"description"`
	NumberOfSub int            `gorm:"column:number_of_sub;not null" json:"number_of_sub"`
	CourseID    uint           `gorm:"column:course_id;not null" json:"course_id"`
	Course      Course         `gorm:"foreignKey:CourseID;references:ID" json:"course"`
	Subchapters []Subchapter   `json:"subchapters"`
}
