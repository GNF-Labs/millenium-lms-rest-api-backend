package endpoints

import (
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/auth"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/handlers"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RegisterAuthRoutes sets up the authentication-related routes
func RegisterAuthRoutes(r *gin.Engine, jwtKey []byte) {
	r.POST("/login", func(c *gin.Context) {
		handlers.HandleLogin(c, jwtKey)
	})

	r.POST("/register", handlers.HandleRegister)

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
}
