package managers

import (
	"errors"
	"fmt"
	"main/common"
	"main/database"
	"main/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type UserManager interface {
	Create(userData *common.UserCreationInput) (*models.User, error)
	List() ([]models.User, error)
	Get(id string) (models.User, error)
	Update(userId string, userData *common.UserUpdationInput) (*models.User, error)
	Delete(id string) (*models.User, error)
	Login(email, password string) (*models.User, string, error)
}

type userManager struct {
	//dbClient
}

func NewUserManager() UserManager {
	return &userManager{}
}

// Create New User
func (userManager *userManager) Create(userData *common.UserCreationInput) (*models.User, error) {

	var existingemail models.User
	value := database.DB.Where("email = ?",userData.Email).First(&existingemail)
	if !errors.Is(value.Error, gorm.ErrRecordNotFound){
		if value.Error == nil {
			return nil,ErrEmailAlreadyExists
		}
		return nil, fmt.Errorf("failed to check email existence %w",value.Error)
	}


	// Generate a UUID for the token
	uuidToken, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate UUID: %w", err)
	}

	newUser := &models.User{
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
		Email:     userData.Email,
		Password:  userData.Password,
		Phone:     userData.Phone,
		Token:     uuidToken.String(), // Set the generated UUID as the token
	}

	//Hash the password
	hashedPassword , err := bcrypt.GenerateFromPassword([]byte(userData.Password),bcrypt.DefaultCost)
	if err != nil{
		return nil,fmt.Errorf("failed to hash the password: %w",err)
	}

	newUser.Password = string(hashedPassword)

	result := database.DB.Create(newUser)
	if result.Error != nil {
		return nil,fmt.Errorf("failed to create a user: %w",result.Error)
	}

	if newUser.ID == 0 {
		
		return nil, errors.New("failed to create a new user")
	}

	return newUser, nil
}



// List All Users
func (userManager *userManager) List() ([]models.User, error) {

	users := []models.User{}
	database.DB.Find(&users)

	return users, nil
}



// Get Single User
func (userManager *userManager) Get(id string) (models.User, error) {

	user := models.User{}
	database.DB.First(&user, id)

	return user, nil
}


//Update the User
func (userManager *userManager) Update(userId string, userData *common.UserUpdationInput) (*models.User, error) {

	user := models.User{}
	result := database.DB.First(&user, userId)
	if result.Error != nil {
		if errors.Is(result.Error, database.DB.Error) {
			return nil, fmt.Errorf("user not found") // Indicate user not found
		}
		return nil, fmt.Errorf("failed to find user: %w", result.Error)
	}

	// Update fields
	user.FirstName = userData.FirstName
	user.LastName = userData.LastName
	user.Email = userData.Email
	user.Phone = userData.Phone

	// Hash password if it's being updated
	if userData.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = string(hashedPassword)
	}

	result = database.DB.Save(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update user: %w", result.Error)
	}

	return &user, nil
}


//Delete the User
func (userManager *userManager) Delete(id string) (*models.User, error) {

	user := &models.User{}
	database.DB.Delete(user, id)

	return user, nil
}


func (userManager *userManager) Login(email, password string) (*models.User, string, error) {

	user := models.User{}
	result := database.DB.Where("email = ?", email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, database.DB.Error) {
			fmt.Println("Login: User not found for email:", email) // Debug log
			return nil, "", errors.New("invalid credentials")
		}
		fmt.Println("Login: Database error:", result.Error) // Debug log
		return nil, "", fmt.Errorf("failed to find user: %w", result.Error)
	}

	fmt.Println("Login: Email:", email)                                   // Debug
	fmt.Println("Login: Stored Hashed Password:", user.Password)            // Debug
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) // Debug
	if err != nil {
		fmt.Println("Login: Password comparison error:", err) // Debug Log
		return nil, "", errors.New("invalid credentials")
	}

	
	token, err := common.GenerateJWT(user.Email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return &user, token, nil
}

