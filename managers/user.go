package managers

import (
	"errors"
	"main/common"
	"main/database"
	"main/models"
)

type UserManager struct {
	//dbClient 
}

func NewUserManager() *UserManager {
	return &UserManager{}
}

func (userManager *UserManager) Create(userData *common.UserCreationInput) (*models.User,error){

	newUser := &models.User{FullName: userData.FullName,Email: userData.Email}
	database.DB.Create(newUser)

	if newUser.ID == 0 {
		return nil,errors.New("failed to create a new user")
	}

	return newUser,nil
}


func (userManager *UserManager) List() ([]models.User,error){

	users := []models.User{}
	database.DB.Find(&users)

	return users, nil
}

