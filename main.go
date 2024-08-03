package main

import (
	"fmt"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/auth"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/databases"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strings"
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

	databaseError := databases.Connect()
	if databaseError != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	fmt.Println("Environment variables loaded successfully")

	// run the gin router context
	r := gin.Default()
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
		if bearerToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		splitToken := strings.Split(bearerToken, " ")
		if len(splitToken) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		tokenString := splitToken[1]
		claims, err := auth.VerifyJWT(jwtKey, tokenString)
		if err != nil {
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
