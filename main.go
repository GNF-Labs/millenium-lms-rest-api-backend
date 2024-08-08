package main

import (
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/auth"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/databases"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/models"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	var jwtKey = []byte(os.Getenv("JWT_SECRET"))
	if len(jwtKey) == 0 {
		log.Fatalf("JWT_SECRET is not set in the environment variables")
	}
	log.Println("Environment variables loaded successfully")
	databaseError := databases.Connect()
	if databaseError != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	log.Println("Database connected successfully")

	// run the gin router context
	r := gin.Default()

	// enable CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))

	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		// Here you would typically check the username and password against your user store
		if username != "testuser" || password != "testpassword" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		token, err := auth.GenerateJWT(jwtKey, username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	})

	r.GET("/hello", func(c *gin.Context) {
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

		c.JSON(http.StatusOK, gin.H{
			"username": claims.Username,
			"token":    claims.RegisteredClaims,
		})
	})

	r.GET("/profile/:username", func(c *gin.Context) {
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
	})

	r.GET("/verify-token", func(c *gin.Context) {
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
		c.JSON(http.StatusOK, gin.H{
			"username": claims.Username,
			"token":    claims.RegisteredClaims,
		})
	})
	err = r.Run("localhost:8080")
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
