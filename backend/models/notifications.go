package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	LimitID   uint    `json:"limit_id"`
	RealValue float64 `json:"real_value"`
	Timestamp time.Time  `json:"timestamp"`

	Limit Limit `gorm:"foreignKey:LimitID"`
}
