package models

import (
	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	LimitID   uint    `json:"limit_id"`
	RealValue float64 `json:"real_value"`
	Timestamp string  `json:"timestamp"`

	Limit Limit `gorm:"foreignKey:LimitID"`
}
