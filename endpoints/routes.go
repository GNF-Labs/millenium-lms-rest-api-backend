package endpoints

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

// RegisterRoutes sets up all routes for the application
func RegisterRoutes(r *gin.Engine, jwtKey []byte) {
	// Enable CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))

	// Register routes
	RegisterAuthRoutes(r, jwtKey)
	RegisterUserRoutes(r, jwtKey)
	RegisterCourseRoutes(r, jwtKey)
}
