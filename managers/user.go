package managers

import (
	"errors"
	"fmt"
	"main/common"
	"main/database"
	"main/models"

	"github.com/google/uuid" // Import the UUID package
)

type UserManager interface {
	Create(userData *common.UserCreationInput) (*models.User, error)
	List() ([]models.User, error)
	Get(id string) (models.User, error)
	Update(userId string, userData *common.UserUpdationInput) (*models.User, error)
	Delete(id string) (*models.User, error)
}

type userManager struct {
	//dbClient
}

func NewUserManager() UserManager {
	return &userManager{}
}

// Create New User
func (userManager *userManager) Create(userData *common.UserCreationInput) (*models.User, error) {
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
	database.DB.Create(newUser)

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

func (userManager *userManager) Update(userId string, userData *common.UserUpdationInput) (*models.User, error) {

	user := models.User{}

	database.DB.First(&user, userId)

	if user.ID == 0 {
		fmt.Println("User is already deleted")
	}

	database.DB.Model(&user).Updates(models.User{
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
		Email:     userData.Email,
		Password:  userData.Password,
		Phone:     userData.Phone,
	})

	return &user, nil
}

func (userManager *userManager) Delete(id string) (*models.User, error) {

	user := &models.User{}
	database.DB.Delete(user, id)

	return user, nil
}

