package models

import (
	"gorm.io/gorm"
	"time"
)

type Course struct {
	gorm.Model
	Name          string     `gorm:"column:name;type:varchar64;not null"`
	Description   string     `gorm:"column:about;type:text;not null"`
	Categories    []Category `gorm:"many2many:course_category;"`
	TimeEstimated time.Time  `gorm:"column:time_estimated;type:datetime;not null"`
	rating        float32    `gorm:"column:rating;type:float;not null"`
}
