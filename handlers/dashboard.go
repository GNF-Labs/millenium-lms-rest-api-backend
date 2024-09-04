package handlers

import (
	"errors"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/auth"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/databases"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type Dashboard struct {
	FullName          string  `json:"name"`
	CompletedCourses  int     `json:"completed_courses"`
	InProgressCourses int     `json:"in_progress_courses"`
	CompletionRate    float64 `json:"completion_rate"`
	LatestCourseName  string  `json:"latest_course_name"`
}

func HandleDashboard(c *gin.Context, jwtKey []byte) {
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
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not allowed to view this user's statistic"})
		return
	}

	var data Dashboard
	if err := dashboardQuery(username).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return empty data if no record is found
			data = Dashboard{}
			c.JSON(http.StatusOK, data)
		} else {
			// Handle other errors
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, data)

}

// dashboardQuery will retrieve all information
// required for dashboard data
func dashboardQuery(username string) *gorm.DB {
	subQuery := databases.DB.Table("user_course_interactions uci_inner").
		Select("uci_inner.course_id").
		Where("uci_inner.user_id = users.id").
		Where("uci_inner.registered = TRUE").
		Order("uci_inner.last_interaction DESC").
		Limit(1)

	return databases.DB.Table("users").
		Select("users.full_name AS name, "+
			"COUNT(CASE WHEN uci.registered = TRUE AND uci.completed = TRUE THEN 1 END) AS completed_courses, "+
			"COUNT(CASE WHEN uci.registered = TRUE AND uci.completed = FALSE THEN 1 END) AS in_progress_courses, "+
			"uci.completion_rate, "+
			"c.name AS latest_course_name").
		Joins("INNER JOIN user_course_interactions uci ON users.id = uci.user_id").
		Joins("LEFT JOIN courses c ON c.id = (?)", subQuery).
		Where("users.username = ?", username).
		Group("users.id, users.full_name, uci.completion_rate, c.name").
		Limit(1)
}
