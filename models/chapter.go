package models

import (
	"encoding/json"
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
	CourseID    uint           `gorm:"primarykey;column:course_id;not null" json:"course_id"`
	Course      Course         `gorm:"foreignKey:CourseID;references:ID" json:"course"`
	Subchapters []Subchapter   `json:"subchapters"`
	Order       uint           `gorm:"column:order;default:0" json:"order"`
}

func (c *Chapter) MarshalJSON() ([]byte, error) {
	type Alias Chapter // Avoid recursion
	return json.Marshal(&struct {
		*Alias
		Course      CourseDTO       `json:"course"`
		Subchapters []SubchapterDTO `json:"subchapters"`
	}{
		Alias:  (*Alias)(c),
		Course: CourseDTO{ID: c.Course.ID, Name: c.Course.Name},
		Subchapters: func() []SubchapterDTO {
			var result []SubchapterDTO
			for _, sub := range c.Subchapters {
				result = append(result, SubchapterDTO{ID: sub.ID, Name: sub.Name})
			}
			return result
		}(),
	})
}
