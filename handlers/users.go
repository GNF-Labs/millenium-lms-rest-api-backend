package handlers

import (
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/auth"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/databases"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/models"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

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
		"user":    userUpdate,
	})
}
