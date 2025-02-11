package main

import (
	"log"
	"os"
	"main/database"
	"main/handlers"
	"main/managers"

	"github.com/gin-gonic/gin"
)


func main() {
	router := gin.Default()

	log.Println("Database Initializing started...")
	database.Initialize()
	log.Println("Database Initializing ended...")

	userManager := managers.NewUserManager()
	userHandler := handlers.NewUserHandlerFrom(userManager)

	userHandler.RegisterUserApis(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	router.Run(":" + port)
}
