package handlers

import (
	"errors"
	"fmt"
	"os"
	"strings"

	// "log" // Import the log package
	"main/common"
	"sync"
	// "main/database"
	"main/managers"
	// "main/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
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

var tokenBlacklist sync.Map

// Middleware to extract user email from JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Authorization header required"})
			return

		}

		tokenString := strings.Replace(authHeader, "Bearer", "", 1)
		tokenString = strings.TrimSpace(tokenString)
		log.Println("Token String:", tokenString) // Log the token after removing "Bearer"

		if _, exists := tokenBlacklist.Load(tokenString); exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token has been logged out"})
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &common.JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if err != nil {
			log.Println("Token Parsing Error:", err) // Log the error
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid Token"})
			return
		}

		claims, ok := token.Claims.(*common.JWTClaim)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token class"})
			return
		}

		c.Set("email", claims.Email)
		c.Next()
	}
}

// Grouping the apis according to the Operations
func (userHandler *UserHandler) RegisterUserApis(router *gin.Engine) {
	userGroup := router.Group(userHandler.groupName)
	userGroup.POST("/signup", userHandler.SignUp)
	userGroup.POST("", userHandler.Create)
	userGroup.GET("", userHandler.List)
	userGroup.GET(":userid/", userHandler.Get)
	userGroup.DELETE(":userid/", userHandler.Delete)
	userGroup.PATCH(":userid/", userHandler.Update)
	userGroup.POST("/login", userHandler.Login)
	userGroup.POST("/logout", userHandler.Logout)
	userGroup.GET("/profile", AuthMiddleware(), userHandler.ViewProfile)

}

func (userHandler *UserHandler) SignUp(ctx *gin.Context) {

	//Calling data
	userData := common.NewUserCreationInput()
	//Bindind Data
	if err := ctx.BindJSON(&userData); err != nil {
		common.BadResponse(ctx, "Invalid signup data")
		return
	}

	newUser, err := userHandler.userManager.Create(userData)
	if err != nil {
		if errors.Is(err, managers.ErrEmailAlreadyExists) {
			common.BadResponse(ctx, "Email Already Exists")
			return
		}

		common.InternalServerErrorResponse(ctx, "Failed to create a message")
		return
	}

	token, err := common.GenerateJWT(newUser.Email)
	if err != nil {
		fmt.Println("Error generating token: ", err)
		common.InternalServerErrorResponse(ctx, "Failed to Generate token")
		return
	}

	common.SuccessResponseWithData(ctx, "Signup Successfull", gin.H{
		"user_id":  newUser.ID,
		"email":    newUser.Email,
		"username": newUser.FirstName,
		"token":    token,
	})

}

// Creating the User Apis
func (userHandler *UserHandler) Create(ctx *gin.Context) {

	userData := common.NewUserCreationInput()

	err := ctx.BindJSON(&userData) //Binding the data

	if err != nil {

		common.BadResponse(ctx, "Failed to bind data")
	}

	newUser, err := userHandler.userManager.Create(userData)

	if err != nil {

		if errors.Is(err, managers.ErrEmailAlreadyExists) {
			common.BadResponse(ctx, "Email already exists")
			return
		}
		common.BadResponse(ctx, "failed to create a user")
	}

	ctx.JSON(http.StatusOK, newUser)
}

// Listing the whole User
func (userHandler *UserHandler) List(ctx *gin.Context) {

	allUser, err := userHandler.userManager.List()
	if err != nil {
		common.BadResponse(ctx, "failed to fetch list")
	}

	ctx.JSON(http.StatusOK, allUser)
}

// Listing the Single User According to the needs.
func (userHandler *UserHandler) Get(ctx *gin.Context) {

	detailUser, ok := ctx.Params.Get("userid")

	if !ok {
		fmt.Println("failed to fetch single user")
	}

	allUser, err := userHandler.userManager.Get(detailUser)

	if allUser.ID == 0 {
		common.BadResponse(ctx, "User is already deleted")

		return
	}

	if err != nil {
		fmt.Println("Failed to fetch the user")
	}

	ctx.JSON(http.StatusOK, allUser)
}

//Update the User
func (userHandler *UserHandler) Update(ctx *gin.Context) {

	userId, ok := ctx.Params.Get("userid")

	if !ok {
		fmt.Println("failed to fetch single user")
	}

	userUpdate := common.NewUserUpdationInput()

	err := ctx.BindJSON(&userUpdate) //Binding the data

	if err != nil {

		common.BadResponse(ctx, "Failed to bind data")
		return
	}

	user, err := userHandler.userManager.Update(userId, userUpdate)

	if err != nil {
		common.BadResponse(ctx, "failed to Update a user")
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// Deleting the user
func (userHandler *UserHandler) Delete(ctx *gin.Context) {

	deleteUser, ok := ctx.Params.Get("userid")

	if !ok {
		fmt.Println("failed to delete user")
	}

	_, err := userHandler.userManager.Delete(deleteUser)

	if err != nil {
		fmt.Println("Failed to delete the user")
	}

	common.SuccessResponse(ctx, "Deleted user")

}

// Login User Function
func (userHandler *UserHandler) Login(ctx *gin.Context) {

	var loginInput common.LoginInput

	if err := ctx.BindJSON(&loginInput); err != nil {
		common.BadResponse(ctx, "Invalid login data")
		return
	}

	user, token, err := userHandler.userManager.Login(loginInput.Email, loginInput.Password)

	if err != nil {
		fmt.Println(err.Error())
		common.BadResponse(ctx, "Invalid Credentials")
		return
	}

	fmt.Println("Login: Token from login API:", token)

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "Login successful",
		"user_id":  user.ID,
		"email":    user.Email,
		"token":    token,
		"username": user.FirstName,
	})
}

// Logout User Function
func (userHandler *UserHandler) Logout(ctx *gin.Context) {
	var logoutInput common.LogoutInput

	if err := ctx.BindJSON(&logoutInput); err != nil {
		common.BadResponse(ctx, "Invalid logout data")
		return
	}

	// Add token to the blacklist:
	tokenBlacklist.Store(logoutInput.Token, true)

	fmt.Println("Logout :Token received by Logout API:", logoutInput.Token)
	err := userHandler.userManager.Logout(logoutInput.Token)

	if err != nil {
		fmt.Println(err.Error())
		common.BadResponse(ctx, "Logout Failed")
		return
	}

	common.SuccessResponse(ctx, "Logout successfull")

}

// View Profile function
func (userHandler *UserHandler) ViewProfile(ctx *gin.Context) {

	email, _ := ctx.Get("email")

	profile, err := userHandler.userManager.ViewProfile(email.(string))

	if err != nil {
		common.InternalServerErrorResponse(ctx, "Failed to fetch profile")
		return
	}

	common.SuccessResponseWithData(ctx, "Profile retrieved successfully", profile)
}
