package handlers

import (
	"errors"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/auth"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/databases"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/models"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"regexp"
	"strconv"
)

func HandleLogin(c *gin.Context, jwtKey []byte) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is empty"})
	}
	// Here you would typically check the username and password against your user store
	var user models.User
	if err := databases.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := auth.GenerateJWT(jwtKey, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": user.Username,
		"token":    token,
	})
}

func HandleProfile(c *gin.Context, jwtKey []byte) {
	username := c.Param("username")
	bearerToken := c.GetHeader("Authorization")
	tokenString, err := utils.ParseToken(bearerToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token not provided"})
		return
	}
	claims, verifyErr := auth.VerifyJWT(jwtKey, tokenString)
	if verifyErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// Verify that the token's username matches the requested profile
	if claims.Username != username {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not allowed to view this profile"})
		return
	}

	var user models.User
	if err := databases.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Return the user profile
	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"username":  user.Username,
		"full_name": user.FullName,
		"email":     user.Email,
		"about":     user.About,
		"role":      user.Role,
		"image_url": user.ImageURL,
	})
}

func HandleVerifyToken(c *gin.Context, jwtKey []byte) {
	username := c.Param("username")
	bearerToken := c.GetHeader("Authorization")
	tokenString, err := utils.ParseToken(bearerToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}
	claims, verifyErr := auth.VerifyJWT(jwtKey, tokenString)
	if verifyErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// Verify that the token's username matches the requested profile
	if claims.Username != username {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not allowed to view this profile"})
		return
	}

	var user models.User
	if err := databases.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Return the user profile
	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"username":  user.Username,
		"full_name": user.FullName,
		"email":     user.Email,
		"about":     user.About,
		"role":      user.Role,
	})
}

func HandleRegister(c *gin.Context) {
	username := c.PostForm("username")
	fullName := c.PostForm("full_name")
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Basic form validation
	if username == "" || fullName == "" || email == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	// Additional validations
	if len(username) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username must be at least 3 characters long"})
		return
	}

	// Validate email using regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	if len(password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters long"})
		return
	}

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}
	user := models.User{
		Username: username,
		FullName: fullName,
		Email:    email,
		Password: string(hashedPassword),
		Role:     "student",
	}
	if err := databases.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

func HandleUpdateProfile(c *gin.Context, jwtKey []byte) {
	username := c.Param("username")
	bearerToken := c.GetHeader("Authorization")
	tokenString, err := utils.ParseToken(bearerToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}
	claims, verifyErr := auth.VerifyJWT(jwtKey, tokenString)
	if verifyErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// Verify that the token's username matches the requested profile
	if claims.Username != username {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not allowed to edit this user"})
		return
	}

	// Fetch the user to be updated
	var user models.User
	if err := databases.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Parse and validate form data
	var userUpdate models.User
	if err := c.ShouldBindJSON(&userUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
		return
	}

	// Update the user's profile
	if err := databases.DB.Model(&models.User{}).Where("username = ?", username).Updates(&userUpdate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "profile updated successfully",
		"user":    user,
	})
}

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
	if err := query.Offset(offset).Limit(pageSize).Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch courses"})
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
	if err := databases.DB.Preload("Chapters").First(&course, "courses.id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"course": course})
}

func GetChapters(c *gin.Context, courseId int) {}
