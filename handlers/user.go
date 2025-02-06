package handlers

import (
	"fmt"
	// "log" // Import the log package
	"main/common"
	// "main/database"
	"main/managers"
	// "main/models"
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

	// var userData struct {
	// 	FullName string `json:"full_name"`
	// 	Email    string `json:"email"`
	// }

	userData := common.NewUserCreationInput()

	err := ctx.BindJSON(&userData)

	fmt.Println(userData.FullName)
	fmt.Println(userData.Email)

	if err!=nil{
		fmt.Println("Failed to bind data")
	}

	// user := &models.User{FullName: userData.FullName, Email: userData.Email}

	// result := database.DB.Create(user) // Store the result
	// if result.Error != nil {
	// 	log.Printf("Database error creating user: %v", result.Error)

	// 	ctx.JSON(http.StatusInternalServerError, gin.H{
	// 		"error":   "Failed to create user",
	// 		"details": result.Error.Error(),
	// 	})
	// }

	newUser, err := userHandler.userManager.Create(userData)
	 
	if err != nil {
		fmt.Println("failed to create a user")
	}

	ctx.JSON(http.StatusOK,newUser)
}
