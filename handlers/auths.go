package handlers

import (
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/auth"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/databases"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/models"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
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
		"user_id":  user.ID,
		"username": user.Username,
		"token":    token,
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
