package main

import (
	"log"
	"main/database"
	"main/handlers"
	"main/managers"
	"os"
	"strconv"

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

	productManager := managers.NewProductManager()
	productHandler := handlers.NewProductHandler(productManager)
	productHandler.RegisterUserApis(router)


	if err := productManager.SeedCategories(); err != nil {
		log.Fatalf("Failed to seed categories: %v", err)
	}
	log.Println("Successfully seeded categories.")

	// Seed products
	seedCountStr := os.Getenv("SEED_PRODUCT_COUNT")
	seedCount := 100
	if seedCountStr != "" {
		newSeedCount, err := strconv.Atoi(seedCountStr)
		if err != nil {
			log.Printf("Invalid SEED_PRODUCT_COUNT value is %s. Using default of 100.", seedCountStr)
			seedCount = 100
		} else {
			seedCount = newSeedCount
		}
	}
	if err := productManager.SeedProducts(seedCount); err != nil {
		log.Fatalf("Failed to seed products %v", err)
	}
	log.Printf("Successfully seeded %d products.",seedCount)


	port := os.Getenv("PORT")
	if port == "" {
		port =  os.Getenv("PORT")
	}
	router.Run(":" + port)
}
