package models

type CourseStructure struct {
	CourseID     int  `gorm:"primaryKey"`
	ChapterID    *int `gorm:"uniqueIndex:idx_course_chapter_subchapter"`
	SubchapterID *int `gorm:"uniqueIndex:idx_course_chapter_subchapter"`

	// Foreign key relation to CourseStructure
	UserProgress UserProgress `gorm:"foreignKey:CourseID,ChapterID,SubchapterID;references:CourseID,ChapterID,SubchapterID;constraint:OnDelete:CASCADE"`
}
