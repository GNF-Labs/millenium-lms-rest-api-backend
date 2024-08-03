package models

import "gorm.io/gorm"

type Chapter struct {
	gorm.Model
	Name        string `gorm:"column:name;type:varchar(50);not null"`
	Description string `gorm:"column:description;type:text"`
	NumberOfSub int    `gorm:"column:number_of_sub;not null"`
	CourseID    uint   `gorm:"column:course_id;not null"`
	Course      Course `gorm:"foreignKey:CourseID;references:ID"`
}
