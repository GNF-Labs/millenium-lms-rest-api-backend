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
	ChapterID uint           `gorm:"column:chapter_id;not null" json:"chapter_id"`
	Chapter   Chapter        `gorm:"foreignKey:ChapterID;references:ID"`
}
