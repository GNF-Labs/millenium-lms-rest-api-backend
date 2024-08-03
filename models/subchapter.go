package models

import "gorm.io/gorm"

type Subchapter struct {
	gorm.Model
	Name      string  `gorm:"column:name;size:64;not null"`
	Content   string  `gorm:"column:content;type:text;not null"`
	ChapterID uint    `gorm:"column:chapter_id;not null"`
	Chapter   Chapter `gorm:"foreignKey:ChapterID;references:ID"`
}
