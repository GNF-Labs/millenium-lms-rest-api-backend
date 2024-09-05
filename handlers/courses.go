package handlers

import (
	"errors"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/databases"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"fmt"
)

type RequestBody struct {
	CoursesID []int `json:"courses_id"`
}

func GetChapterDetail(c *gin.Context) {
	var err error
	courseId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chapterId, err := strconv.Atoi(c.Param("chapter_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var chapter models.Chapter
	if err = databases.DB.
	Preload("Subchapters", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name, chapter_id") // Select only necessary columns for subchapters
	}).Preload("Course", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name") // Select only necessary columns for subchapters
	}).
		Joins("JOIN courses ON courses.id = chapters.course_id").
		Joins("LEFT JOIN subchapters ON subchapters.chapter_id = chapters.id").
		Where("chapters.course_id = ? AND chapters.id = ?", courseId, chapterId).
		Select("chapters.id, chapters.name, chapters.description, chapters.number_of_sub, chapters.order, courses.id as course_id, courses.name as course_name, subchapters.id as subchapter_id, subchapters.name as subchapter_name").
		First(&chapter).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Chapter not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}


	c.JSON(http.StatusOK, gin.H{"message": "Data retrieved", "data": chapter})
}

func GetSubchaptersFromChapter(c *gin.Context) {
	courseIdStr := c.Param("id")
	chapterIdStr := c.Param("chapter_id")
	subchapterIdStr := c.Param("subchapter_id")

	fmt.Println(courseIdStr)
	// courseId, err := strconv.Atoi(courseIdStr)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course id"})
	// 	return
	// }

	chapterId, err := strconv.Atoi(chapterIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chapter id"})
		return
	}

	subchapterId, err := strconv.Atoi(subchapterIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subchapter id"})
		return
	}

	var subchapter models.Subchapter
	if err = databases.DB.
		Joins("JOIN chapters ON chapters.id = subchapters.chapter_id").
		Where("subchapters.chapter_id = ? AND subchapters.id = ?", chapterId, subchapterId).
		Select("subchapters.id, subchapters.name, chapters.id AS chapter_id, chapters.name AS chapter_name").
		Joins("JOIN courses ON courses.id = chapters.course_id").
		Select("courses.id AS course_id, courses.name AS course_name, chapters.id as chapter_id, chapters.name as chapter_name, subchapters.id, subchapters.name, subchapters.content, subchapters.order").
		First(&subchapter).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subchapter not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data retrieved", "data": subchapter})

}

func GetCoursesByIdCollection(c *gin.Context) {
	var requestBody RequestBody

	// Bind the request body to the RequestBody struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if courses_id is provided
	if len(requestBody.CoursesID) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No course IDs provided"})
		return
	}

	// Fetch courses from the database
	var courses []models.Course
	if err := databases.DB.
		Where("id IN ?", requestBody.CoursesID).
		Select("id", "name").
		Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with the retrieved courses
	c.JSON(http.StatusOK, gin.H{"message": "Data retrieved", "data": courses})
}
