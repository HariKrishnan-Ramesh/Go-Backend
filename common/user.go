package common

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"time"
)

type Address struct {
	Address1 string `json:"address1"`
	Address2 string `json:"address2,omitempty"`
	City     string `json:"city"`
	District string `json:"district,omitempty"`
	State    string `json:"state"`
	Country  string `json:"country"`
	Pin      string `json:"pin"`
}

type UserCreationInput struct {
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Email     string  `json:"email" binding:"required"`
	Phone     string  `json:"phone"`
	Password  string  `json:"password" binding:"required"`
	Token     string  `gorm:"uniqueIndex" json:"token"`
	Address   Address `json:"address"`
	Image     string  `json:"image,omitempty"`
}

type UserUpdationInput struct {
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Email     string  `json:"email"`
	Phone     string  `json:"phone"`
	Password  string  `json:"password"`
	Token     string  `gorm:"uniqueIndex" json:"token"`
	Address   Address `json:"address"`
	Image     string  `json:"image,omitempty"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

type LogoutInput struct {
	Token string `json:"token" binding:"required"`
}

type requestResponse struct {
	Message string      `json:"message"`
	Status  uint        `json:"status"`
	Data    interface{} `json:"data,omitempty"`
}

type ProfileResponse struct {
	Id        uint    `json:"id"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Email     string  `json:"email"`
	Phone     string  `json:"phone"`
	Address   Address `json:"address"`
	Image     string  `json:"image,omitempty"`
}

func NewUserCreationInput() *UserCreationInput {
	return &UserCreationInput{}
}

func NewUserUpdationInput() *UserUpdationInput {
	return &UserUpdationInput{}
}

func SuccessResponse(ctx *gin.Context, msg string) {
	response := requestResponse{
		Message: msg,
		Status:  http.StatusOK,
	}

	ctx.JSON(http.StatusOK, response)
}

func SuccessResponseWithData(ctx *gin.Context, msg string, data interface{}) {
	response := requestResponse{
		Message: msg,
		Status:  http.StatusOK,
		Data:    data,
	}
	ctx.JSON(http.StatusOK, response)
}

func BadResponse(ctx *gin.Context, msg string) {
	response := requestResponse{
		Message: msg,
		Status:  http.StatusBadRequest,
	}

	ctx.JSON(http.StatusBadRequest, response)
}

func InternalServerErrorResponse(ctx *gin.Context, msg string) {
	response := requestResponse{
		Message: msg,
		Status:  http.StatusInternalServerError,
	}
	ctx.JSON(http.StatusInternalServerError, response)
}

func GenerateJWT(email string) (string, error) {
	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &JWTClaim{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type JWTClaim struct {
	Email string
	jwt.RegisteredClaims
}

