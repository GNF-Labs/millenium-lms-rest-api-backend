package models

type CourseCategory struct {
	CourseID   uint `gorm:"column:course_id;primaryKey"`
	CategoryID uint `gorm:"column:category_id;primaryKey"`
}
