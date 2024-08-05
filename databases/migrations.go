package databases

import (
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := Connect()
	if err != nil {
		return err
	}
	migrateErr := db.AutoMigrate(&models.User{},
		&models.Course{},
		&models.Chapter{},
		&models.Category{},
		&models.Review{},
		&models.Subchapter{},
		&models.UserCourseInteraction{},
	)
	if migrateErr != nil {
		return migrateErr
	}
	return nil
}
