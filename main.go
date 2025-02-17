package main

import (
	"log"
	"os"
	"main/database"
	"main/handlers"
	"main/managers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)


func main() {
	router := gin.Default()

	err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

	log.Println("Database Initializing started...")
	database.Initialize()
	log.Println("Database Initializing ended...")

	userManager := managers.NewUserManager()
	userHandler := handlers.NewUserHandlerFrom(userManager)

	userHandler.RegisterUserApis(router)

	port := os.Getenv("PORT")
	if port == "" {
		port =  os.Getenv("PORT")
	}
	router.Run(":" + port)
}
