package endpoints

import (
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/handlers"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// RegisterCourseRoutes sets up the course-related routes
func RegisterCourseRoutes(r *gin.Engine, jwtKey []byte) {
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

	r.GET("/dashboard/:username", func(c *gin.Context) {
		handlers.HandleDashboard(c, jwtKey)
	})

	r.GET("/courses/:id/:chapter_id", handlers.GetChapterDetail)

	r.GET("/courses/:id/:chapter_id/:subchapter_id", func(context *gin.Context) {
		handlers.GetSubchaptersFromChapter(context)
	})
}
