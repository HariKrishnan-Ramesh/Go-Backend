package database

import (
	"fmt"
	"log"
	"main/models"
	"os"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Initialize() {
	var err error

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	//log.Fatalf("Connection established %w",dbUser) Like this {other functions also printed} can be used to DEBUG to know the connection established is correct or not.

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
	log.Printf("DSN: %s", dsn) //This line is for debugging


	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
		panic("Failed to connect database")
	}

	err = DB.AutoMigrate(&models.User{},&models.Category{},&models.Product{}, &models.Wishlist{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
		panic("Failed to Automigrate database")
	}

	log.Println("Database connection established and auto-migration complete.")

}
