/* migrate.go */

// Package databases is package that handle database things
package main

import (
	"flag"
	"fmt"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/databases"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	host := flag.String("host", "localhost", "Database host")
	port := flag.Int("port", 5432, "Database port")
	user := flag.String("user", "postgres", "Database user")
	password := flag.String("password", "", "Database password")
	dbname := flag.String("dbname", "mydb", "Database name")

	// Parse command-line arguments
	flag.Parse()

	// Create a connection string
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta",
		*host, *port, *user, *password, *dbname)

	log.Println("Try to connect to the database...")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Successfully connected to database")
	err = databases.Migrate(db)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	fmt.Println("Migration complete")
}
