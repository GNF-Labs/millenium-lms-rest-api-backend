package main

import (
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/api"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/auth"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/databases"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/handlers"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	var err error
	err = godotenv.Load()
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

	err = api.InitGoogleStorageClient()
	if err != nil {
		log.Fatalf("Could not connect to Google: %v", err)
	}

	log.Println("Google Storage Client Initialized")

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

	// Login Endpoint
	r.POST("/login", func(c *gin.Context) {
		handlers.HandleLogin(c, jwtKey)
	})

	r.POST("/register", handlers.HandleRegister)

	// Hello (for testing)
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
		handlers.HandleProfile(c, jwtKey)
	})

	r.PUT("/profile/:username", func(c *gin.Context) {
		handlers.HandleUpdateProfile(c, jwtKey)
	})

	r.GET("/verify-token/:username", func(c *gin.Context) {
		handlers.HandleVerifyToken(c, jwtKey)
	})

	r.PUT("/interact", func(c *gin.Context) {
		handlers.HandleUserCourseInteractions(c, jwtKey)
	})

	// course endpoints
	r.GET("/courses", func(c *gin.Context) {
		pageStr := c.DefaultQuery("page", "1")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
			return
		}
		searchQuery := c.DefaultQuery("q", "")
		handlers.GetCourses(c, page, searchQuery)
	})
	r.GET("/courses/:id", handlers.GetCourseById)

	// run the server
	err = r.Run("localhost:8080")
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
