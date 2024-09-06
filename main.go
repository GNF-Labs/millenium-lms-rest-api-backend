package main

import (
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/databases"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/endpoints"
	"github.com/GNF-Labs/millenium-lms-rest-api-backend/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
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

	err = services.InitGoogleStorageClient()
	if err != nil {
		log.Fatalf("Could not connect to Google: %v", err)
	}

	log.Println("Google Storage Client Initialized")

	// run the gin router context
	r := gin.Default()

	// enable CORS
	endpoints.RegisterRoutes(r, jwtKey)

	// run the server
	err = r.Run("localhost:8080")
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
