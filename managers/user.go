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


//Create New User
func (userManager *UserManager) Create(userData *common.UserCreationInput) (*models.User,error){

	newUser := &models.User{FullName: userData.FullName,Email: userData.Email}
	database.DB.Create(newUser)

	if newUser.ID == 0 {
		return nil,errors.New("failed to create a new user")
	}

	return newUser,nil
}


//List All Users
func (userManager *UserManager) List() ([]models.User,error){

	users := []models.User{}
	database.DB.Find(&users)

	return users, nil
}


//Get Single User
func (userManager *UserManager) Get(id string) (models.User,error){

	user := models.User{}
	database.DB.First(&user, id)

	return user, nil
}
