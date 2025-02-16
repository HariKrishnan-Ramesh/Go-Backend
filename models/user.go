package models

import(


	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string         `json:"firstName"`
	LastName  string         `json:"lastName"` 
	Email     string         `json:"email"`
	Password  string         `json:"password"`  
	Phone     string         `json:"phone"`   
	Token     string         `gorm:"uniqueIndex:idx_users_token,length:191" json:"-"` 
}