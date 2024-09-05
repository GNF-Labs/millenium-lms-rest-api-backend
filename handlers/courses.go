package handlers

import (
	"errors"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/databases"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

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
		Where("course_id = ? AND id = ?", courseId, chapterId).Preload("Course", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	}).Preload("Subchapters", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	}).
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

	courseId, err := strconv.Atoi(courseIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course id"})
		return
	}

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
		Where("course_id = ? AND chapter_id = ? AND id = ?", courseId, chapterId, subchapterId).Preload("Course", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	}).
		Preload("Chapter", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
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
