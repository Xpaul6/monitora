package models

import (
	"gorm.io/gorm"
)

type Component struct {
	gorm.Model
	ServerID uint   `json:"server_id"`
	Type     string `json:"type"`
	Address  string `json:"address"`

	Server Server `gorm:"foreignKey:ServerID"`
	RawLogs []RawLog `gorm:"foreignKey:ComponentID"`
	Limits  []Limit  `gorm:"foreignKey:ComponentID"`
}
