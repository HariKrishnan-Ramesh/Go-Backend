package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserCreationInput struct {
	FullName string `json:"full_name"`
	// LastName          string `json:"last_name"`
	Email string `json:"email"`
	// Phone             string `json:"phone"`
	// Password          string `json:"password"`
	// RegistrationToken string `json:"registration_token"`
}

func NewUserCreationInput() *UserCreationInput {
	return &UserCreationInput{}
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
