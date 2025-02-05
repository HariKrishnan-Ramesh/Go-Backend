package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FullName string `gorm:"type:varchar(100);not null"`
	Email    string `gorm:"type:varchar(100);uniqueIndex;not null"`
}