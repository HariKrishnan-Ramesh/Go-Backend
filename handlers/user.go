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

//Grouping the apis according to the Operations
func (userHandler *UserHandler) RegisterUserApis(router *gin.Engine) {
	userGroup := router.Group(userHandler.groupName)
	userGroup.POST("", userHandler.Create)
	userGroup.GET("",userHandler.List)
	userGroup.GET(":userid/",userHandler.Get)
}


//Creating the User Apis
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


//Listing the whole User
func (userHandler *UserHandler) List(ctx *gin.Context) {

	allUser, err := userHandler.userManager.List()
	if err != nil {
		fmt.Println("failed to fetch list")
	}

	ctx.JSON(http.StatusOK,allUser)
}


//Listing the Single User According to the needs.
func (userHandler *UserHandler) Get(ctx *gin.Context) {

	detailUser , ok := ctx.Params.Get("userid")
	 
	if !ok{
		fmt.Println("failed to fetch single user")
	}

	allUser , err := userHandler.userManager.Get(detailUser)
	if err!=nil{
		fmt.Println("Failed to fetch the user")
	}
	
	ctx.JSON(http.StatusOK,allUser)
}

