package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/GNF-Labs/millenium-lms-rest-api-backend/databases"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
			return db.Select("id, name, chapter_id, subchapters.order") // Select only necessary columns for subchapters
		}).Preload("Course", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name") // Select only necessary columns for subchapters
	}).
		Joins("JOIN courses ON courses.id = chapters.course_id").
		Joins("LEFT JOIN subchapters ON subchapters.chapter_id = chapters.id").
		Where("chapters.course_id = ? AND chapters.id = ?", courseId, chapterId).
		Select("chapters.id, chapters.name, chapters.description, chapters.number_of_sub, chapters.order, courses.id as course_id, courses.name as course_name, subchapters.id as subchapter_id, subchapters.name as subchapter_name, subchapters.order").
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

type SubchapterResponse struct {
	models.Subchapter
	ChapterID   uint   `json:"chapter_id"`
	ChapterName string `json:"chapter_name"`
	CourseID    uint   `json:"course_id"`
	CourseName  string `json:"course_name"`
}

func GetSubchaptersFromChapter(c *gin.Context) {
	courseIdStr := c.Param("id")
	chapterIdStr := c.Param("chapter_id")
	subchapterIdStr := c.Param("subchapter_id")

	fmt.Println(courseIdStr)
	courseId, err := strconv.Atoi(courseIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course id"})
		return
	}

	chapterId, err := strconv.Atoi(chapterIdStr)
	fmt.Println(chapterId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chapter id"})
		return
	}

	subchapterId, err := strconv.Atoi(subchapterIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subchapter id"})
		return
	}

	var result SubchapterResponse
	err = databases.DB.Table("subchapters").
		Select("subchapters.*, chapters.id as chapter_id, chapters.name as chapter_name, courses.id as course_id, courses.name as course_name").
		Joins("JOIN chapters ON chapters.id = subchapters.chapter_id").
		Joins("JOIN courses ON courses.id = chapters.course_id").
		Where("subchapters.id = ? AND chapters.id = ? AND courses.id = ?", subchapterId, chapterId, courseId).
		First(&result).Error

	c.JSON(http.StatusOK, gin.H{"message": "Data retrieved", "data": result})

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
		Select("id", "name", "image_url", "rating", "time_estimated").
		Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with the retrieved courses
	c.JSON(http.StatusOK, gin.H{"message": "Data retrieved", "data": courses})
}

func SetCompleted(c *gin.Context) {
	userIdStr := c.Param("id")
	if userIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No user id provided"})
		return
	}
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}
	var requestBody struct {
		CourseID     int `json:"course_id"`
		ChapterID    int `json:"chapter_id"`
		SubchapterID int `json:"subchapter_id"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
	}
	subChapterId := requestBody.SubchapterID
	chapterId := requestBody.ChapterID
	courseId := requestBody.CourseID

	// update in user progress
	var userProgress models.UserProgress
	err = databases.DB.Where("user_id = ? AND course_id = ? AND chapter_id = ? AND subchapter_id = ?",
		userId, courseId, chapterId, subChapterId).First(&userProgress).Error

	// If the record is found, update the "completed" field
	if err == nil {
		err = databases.DB.Model(&userProgress).Update("completed", true).Error
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User progress not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully Updated the data", "data": userProgress})
}

func CreateUserProgress(c *gin.Context, jwtKey []byte) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide valid username"})
	}
	var requestBody struct {
		CourseID     int `json:"course_id"`
		ChapterID    int `json:"chapter_id"`
		SubchapterID int `json:"subchapter_id"`
	}
	var err error
	if err = c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, ok := CheckPermission(c, jwtKey, username)
	if !ok {
		return
	}

	var createdUserProgress models.UserProgress
	err = databases.DB.Where("user_id = ? AND course_id = ? AND chapter_id = ? AND subchapter_id = ?",
		user.ID, requestBody.CourseID, requestBody.ChapterID, requestBody.SubchapterID).Attrs(&models.UserProgress{
		UserID:       int(user.ID),
		CourseID:     requestBody.CourseID,
		ChapterID:    requestBody.ChapterID,
		SubchapterID: requestBody.SubchapterID,
	}).
		FirstOrCreate(&createdUserProgress).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User progress not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully Created the data", "data": createdUserProgress})
}
