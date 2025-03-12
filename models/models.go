package models

import (
	"gorm.io/gorm"
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

type User struct {
	//gorm.Model
	Id                uint           `gorm:"primarykey" json:"id"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         time.Time      `json:"updatedAt"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
	FirstName         string         `json:"firstName"`
	LastName          string         `json:"lastName"`
	Email             string         `json:"email"`
	Password          string         `json:"password"`
	Phone             string         `json:"phone"`
	VerificationToken string         `gorm:"column:verification_token" json:"-"`
	IsVerified        bool           `gorm:"column:is_verified;default:false" json:"isVerified"`
	Token             string         `gorm:"uniqueIndex:idx_users_token,length:191" json:"-"`
	Address           Address        `json:"address" gorm:"embedded"`
	Image             string         `json:"image,omitempty"`
}

type Product struct {
	//gorm.Model
	Id          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	SKU         string         `json:"sku" gorm:"uniqueIndex:idx_products_sku,length:191"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Price       string         `json:"price"`
	Image       string         `json:"image,omitempty"`
	CategoryID  uint           `json:"categoryID"`
	Category    Category       `json:"category" gorm:"foreignKey:CategoryID"`
}

type Category struct {
	Id          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	ParentID    *uint          `json:"parentID"`
	Children    []Category     `json:"children" gorm:"foreignKey:ParentID"`
	Products    []Product      `json:"products" gorm:"foreignKey:CategoryID"`
}

type Wishlist struct {
	Id        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	UserID    uint      `json:"userID"`
	ProductID uint      `json:"productID"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID"`
}

type Cart struct {
	Id        uint      `gorm:"primaryKey" json:"Id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	UserID    uint      `json:"userID"`
	ProductID uint      `json:"productID"`
	Quantity  uint      `json:"quantity"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID"`
}

type Otp struct{
	ID uint `gorm:"primaryKey"`
	UserID uint `json:"userId"`
	User User `gorm:"foreignKey:UserID"`
	OTP string 	`json:"otp"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}