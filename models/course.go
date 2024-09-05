package models

import (
	"gorm.io/gorm"
	"time"
)

type CourseDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ChapterDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
type SubchapterDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Course struct {
	ID            uint `gorm:"primarykey;auto_increment" json:"id"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	Name          string         `gorm:"column:name;type:varchar(64);not null" json:"name"`
	Description   string         `gorm:"column:description;type:text;not null;default:''" json:"description"`
	TimeEstimated uint           `gorm:"column:time_estimated;not null" json:"time_estimated"`
	Rating        float32        `gorm:"column:rating;type:float;not null;" json:"rating"`
	CategoryID    uint           `gorm:"column:category_id" json:"category_id"`
	Chapters      []Chapter      `json:"chapters"`

	Category Category `gorm:"foreignKey:CategoryID;references:ID"`
}
