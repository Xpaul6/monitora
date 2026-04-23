package models

import (
	"gorm.io/gorm"
)

type MetricType struct {
	gorm.Model
	Name        string `json:"name"`
	Unit        string `json:"unit"`
	Description string `json:"desription"`

	RawLogs []RawLog `gorm:"foreignKey:MetricTypeID"`
	Limits  []Limit  `gorm:"foreignKey:MetricTypeID"`
}
