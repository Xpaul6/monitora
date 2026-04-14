package models

import (
	"gorm.io/gorm"
)

type Server struct {
	gorm.Model
	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
	IP     string `json:"ip"`
	Status string `json:"status"`

	User User `gorm:"foreignKey:UserID"`
	Components []Component `gorm:"foreignKey:ServerID"`
}
