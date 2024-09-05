package endpoints

import (
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/handlers"
	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes sets up the user-related routes
func RegisterUserRoutes(r *gin.Engine, jwtKey []byte) {
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

	r.GET("/interact/:username", func(c *gin.Context) {
		handlers.GetUserCourseInteractions(c, jwtKey)
	})
}
