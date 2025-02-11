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
	userManager managers.UserManager
}

func NewUserHandlerFrom(userManager managers.UserManager) *UserHandler {
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
	userGroup.DELETE(":userid/",userHandler.Delete)
	userGroup.PATCH(":userid/",userHandler.Update)
	userGroup.POST("/login",userHandler.Login)
}


//Creating the User Apis
func (userHandler *UserHandler) Create(ctx *gin.Context) {

	userData := common.NewUserCreationInput()

	err := ctx.BindJSON(&userData) //Binding the data

	if err!=nil{

		common.BadResponse(ctx, "Failed to bind data")
	}

	newUser, err := userHandler.userManager.Create(userData)
	 
	if err != nil {
		common.BadResponse(ctx, "failed to create a user")
	}

	ctx.JSON(http.StatusOK,newUser)
}


//Listing the whole User
func (userHandler *UserHandler) List(ctx *gin.Context) {

	allUser, err := userHandler.userManager.List()
	if err != nil {
		common.BadResponse(ctx, "failed to fetch list")
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

	if allUser.ID == 0 {
		common.BadResponse(ctx, "User is already deleted")

		return
	}


	if err!=nil{
		fmt.Println("Failed to fetch the user")
	}

	ctx.JSON(http.StatusOK,allUser)
}



//Update the User

func (userHandler *UserHandler) Update(ctx *gin.Context) {

	userId , ok := ctx.Params.Get("userid")
	 
	if !ok{
		fmt.Println("failed to fetch single user")
	}

	userUpdate := common.NewUserUpdationInput()

	err := ctx.BindJSON(&userUpdate) //Binding the data

	if err!=nil{
		
		common.BadResponse(ctx, "Failed to bind data")
		return
	}


	user, err := userHandler.userManager.Update(userId, userUpdate)
	 
	if err != nil {
		common.BadResponse(ctx, "failed to Update a user")
		return
	}

	ctx.JSON(http.StatusOK,user)
}



//Deleting the user
func (userHandler *UserHandler) Delete(ctx *gin.Context) {

	deleteUser , ok := ctx.Params.Get("userid")
	 
	if !ok{
		fmt.Println("failed to delete user")
	}

	_,err := userHandler.userManager.Delete(deleteUser)

	if err!=nil{
		fmt.Println("Failed to delete the user")
	}
	
	common.SuccessResponse(ctx , "Deleted user")

}


func (userHandler *UserHandler) Login(ctx *gin.Context) {

	var loginInput common.LoginInput

	if err := ctx.BindJSON(&loginInput) ; err != nil {
		common.BadResponse(ctx, "Invalid login data")
		return
	}

	user, token , err := userHandler.userManager.Login(loginInput.Email,loginInput.Password)

	if err != nil {
		fmt.Println(err.Error())
		common.BadResponse(ctx, "Invalid Credentials")
		return
	}

	ctx.JSON(http.StatusOK,gin.H{
		"message":  "Login successful",
		"user_id":  user.ID,
		"email":    user.Email,
		"token":    token,
		"username": user.FirstName,
	})
}
