package models

import (
	"gorm.io/gorm"
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

type User struct {
	gorm.Model
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	Phone     string  `json:"phone"`
	Token     string  `gorm:"uniqueIndex:idx_users_token,length:191" json:"-"`
	Address   Address `json:"address" gorm:"embedded"`
	Image     string  `json:"image,omitempty"`
}

type Product struct {
	gorm.Model
	SKU         string `json:"sku" gorm:"uniqueIndex:idx_products_sku,length:191"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Image       string `json:"image,omitempty"`
	CategoryID  uint   `json:"categoryID"`
}
