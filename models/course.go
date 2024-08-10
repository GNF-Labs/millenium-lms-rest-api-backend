package models

import (
	"gorm.io/gorm"
	"time"
)

type Course struct {
	ID            uint `gorm:"primarykey;auto_increment" json:"id"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	Name          string         `gorm:"column:name;type:varchar(64);not null" json:"name"`
	Description   string         `gorm:"column:description;type:text;not null;default:''" json:"description"`
	Categories    []Category     `gorm:"many2many:course_category" json:"categories"`
	TimeEstimated uint           `gorm:"column:time_estimated;not null" json:"time_estimated"`
	Rating        float32        `gorm:"column:rating;type:float;not null;" json:"rating"`
}
