package main

import (
	"log"
	"main/database"
	"main/handlers"
	"main/managers"

	"github.com/gin-gonic/gin"
)


func main() {
	router := gin.Default()

	log.Println("Database Initializing started...")
	database.Intialize()
	log.Println("Database Initializing ended...")


	userManager := managers.NewUserManager()
	userHandler := handlers.NewUserHandlerFrom(userManager)

	userHandler.RegisterUserApis(router)


	
	router.Run() // listen and serve on 0.0.0.0:8080
}
