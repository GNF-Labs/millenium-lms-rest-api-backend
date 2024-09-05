package handlers

import (
	"errors"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/auth"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/databases"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/models"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func HandleUserCourseInteractions(c *gin.Context, jwtKey []byte) {
	// Parse the JSON body
	var requestData models.UserCourseInteraction
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
		return
	}

	// Extract and verify the JWT token
	bearerToken := c.GetHeader("Authorization")
	tokenString, err := utils.ParseToken(bearerToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	claims, verifyErr := auth.VerifyJWT(jwtKey, tokenString)
	if verifyErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// Verify that the token's username matches the user ID in the request data
	var user models.User
	if err := databases.DB.Where("id = ?", requestData.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if claims.Username != user.Username {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not allowed to modify this data"})
		return
	}

	// Upsert operation
	var existingInteraction models.UserCourseInteraction
	if err := databases.DB.Where("user_id = ? AND course_id = ?", requestData.UserID, requestData.CourseID).First(&existingInteraction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Record doesn't exist, create a new one
			if err := databases.DB.Create(&requestData).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create interaction"})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"message": "interaction created successfully", "data": requestData})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch interaction"})
		}
		return
	}

	// Update the existing interaction with the provided data
	if err := databases.DB.Model(&existingInteraction).Updates(requestData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update interaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "interaction updated successfully", "data": existingInteraction})
}

func GetUserCourseInteractions(c *gin.Context, jwtKey []byte) {
	var err error
	username := c.Param("username")
	user, ok := CheckPermission(c, jwtKey, username)
	if !ok {
		return
	}

	var interactions []models.UserCourseInteraction
	if err = databases.DB.Where("user_id = ?", user.ID).Find(&interactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch interactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"interactions": interactions})
}

func GetCourses(c *gin.Context, page int, searchQuery string) {
	var courses []models.Course
	var totalRecords int64
	pageSize := 10 // Define the number of records per page

	// Apply the search filter if a search query is provided
	query := databases.DB.Model(&models.Course{})
	if searchQuery != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	// Get the total count of records that match the query
	if err := query.Count(&totalRecords).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count courses"})
		return
	}

	// Calculate the offset for pagination
	offset := (page - 1) * pageSize

	// Retrieve the records based on the pagination and search query
	if err := query.Offset(offset).Limit(pageSize).Preload("Category").Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Return the results along with pagination metadata
	c.JSON(http.StatusOK, gin.H{
		"total_records": totalRecords,
		"page":          page,
		"page_size":     pageSize,
		"total_pages":   (totalRecords + int64(pageSize) - 1) / int64(pageSize),
		"courses":       courses,
	})

}

func GetCourseById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course id"})
		return
	}
	var course models.Course
	if err := databases.DB.Preload("Chapters", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	}).Preload("Category").First(&course, "courses.id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"course": course})
}
