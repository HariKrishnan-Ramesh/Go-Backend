package managers

import "main/models"

type UserManager struct {
	//dbClient 
}

func NewUserManager() *UserManager {
	return &UserManager{}
}

func (userManager *UserManager) Create(user *models.User) (*models.User,error){
	return nil,nil
}