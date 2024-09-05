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
	migrateErr := db.AutoMigrate(
		&models.User{},
		&models.Course{},
		&models.Chapter{},
		&models.Category{},
		&models.Review{},
		&models.Subchapter{},
		&models.UserCourseInteraction{},
		&models.CourseStructure{},
		&models.UserProgress{},
	)
	if migrateErr != nil {
		return migrateErr
	}

	//modelsToCheck := []interface{}{
	//	&models.User{},
	//	&models.Course{},
	//	&models.Chapter{},
	//	&models.Category{},
	//	&models.Review{},
	//	&models.Subchapter{},
	//	&models.UserCourseInteraction{},
	//}

	//for _, model := range modelsToCheck {
	//	err = dropUnusedColumns(db, model)
	//	if err != nil {
	//		return err
	//	}
	//}

	return nil
}

//func dropUnusedColumns(db *gorm.DB, model interface{}) error {
//	// Fetch the current table structure
//	stmt := &gorm.Statement{DB: db}
//	if err := stmt.Parse(model); err != nil {
//		return err
//	}
//
//	tableName := stmt.Schema.Table
//
//	// Get the current columns in the table
//	var currentColumns []string
//	rows, err := db.Raw("SELECT column_name FROM information_schema.columns WHERE table_name = ?", tableName).Rows()
//	if err != nil {
//		return err
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		var columnName string
//		if err := rows.Scan(&columnName); err != nil {
//			return err
//		}
//		currentColumns = append(currentColumns, columnName)
//	}
//
//	// Check for columns in the table that are not in the model
//	for _, dbColumn := range currentColumns {
//		if _, ok := stmt.Schema.FieldsByDBName[dbColumn]; !ok {
//			log.Printf("Dropping unused column %s from table %s", dbColumn, tableName)
//			if err := db.Migrator().DropColumn(model, dbColumn); err != nil {
//				return err
//			}
//		}
//	}
//
//	return nil
//}
