package handlers

import (
	"log" // Import the log package
	"main/database"
	"main/managers"
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	groupName   string
	userManager *managers.UserManager
}

func NewUserHandlerFrom(userManager *managers.UserManager) *UserHandler {
	return &UserHandler{
		"api/user",
		userManager,
	}
}

func (userHandler *UserHandler) RegisterUserApis(router *gin.Engine) {
	userGroup := router.Group(userHandler.groupName)
	userGroup.POST("", userHandler.Create)
}

func (userHandler *UserHandler) Create(ctx *gin.Context) {
	
	user := &models.User{FullName: "Steven", Email: "steven@1234.com"}

	result := database.DB.Create(user) // Store the result
	if result.Error != nil {
		log.Printf("Database error creating user: %v", result.Error)

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user", 
			"details": result.Error.Error(),  
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
		"userID":  user.ID,
	})
}