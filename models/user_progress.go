package models

type UserProgress struct {
	UserID       int  `gorm:"primaryKey"`
	CourseID     int  `gorm:"primaryKey"`
	ChapterID    int  `gorm:"primaryKey"`
	SubchapterID int  `gorm:"primaryKey"`
	Completed    bool `gorm:"default:false"`
}
