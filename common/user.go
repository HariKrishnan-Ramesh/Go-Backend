package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserCreationInput struct {
	
	FirstName string         `json:"firstName"`
	LastName  string         `json:"lastName"` 
	Email     string         `json:"email"`
	Phone     string          `json:"phone"`   
	Password  string         `json:"password"`  
	Token     string         `gorm:"uniqueIndex" json:"token"` 
}

type UserUpdationInput struct {

	FirstName string         `json:"firstName"`
	LastName  string         `json:"lastName"` 
	Email     string         `json:"email"`
	Phone     string          `json:"phone"`   
	Password  string         `json:"password"`  
	Token     string         `gorm:"uniqueIndex" json:"token"` 
}

func NewUserCreationInput() *UserCreationInput {
	return &UserCreationInput{}
}

func NewUserUpdationInput() *UserUpdationInput {
	return &UserUpdationInput{}
}

type requestResponse struct {
	Message string `json:"message"`
	Status uint `json:"status"`
}

func SuccessResponse(ctx *gin.Context , msg string) {
	response := requestResponse{
		Message: msg,           
		Status:  http.StatusOK, 
	}

	ctx.JSON(http.StatusOK, response) 
}

func BadResponse(ctx *gin.Context,msg string){
	response := requestResponse{
		Message: msg,           
		Status:  http.StatusBadRequest, 
	}

	ctx.JSON(http.StatusBadRequest, response) 
}
