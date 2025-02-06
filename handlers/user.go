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
	userGroup.GET("",userHandler.List)
}

func (userHandler *UserHandler) Create(ctx *gin.Context) {

	userData := common.NewUserCreationInput()

	err := ctx.BindJSON(&userData) //Binding the data

	fmt.Println(userData.FullName)
	fmt.Println(userData.Email)

	if err!=nil{
		fmt.Println("Failed to bind data")
	}

	newUser, err := userHandler.userManager.Create(userData)
	 
	if err != nil {
		fmt.Println("failed to create a user")
	}

	ctx.JSON(http.StatusOK,newUser)
}


func (userHandler *UserHandler) List(ctx *gin.Context) {

	allUser, err := userHandler.userManager.List()
	 
	if err != nil {
		fmt.Println("failed to fetch list")
	}

	ctx.JSON(http.StatusOK,allUser)
}


