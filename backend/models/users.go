package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`

	Servers []Server `gorm:"foreignKey:UserID"`
}
