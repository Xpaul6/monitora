package models

import (
	"gorm.io/gorm"
)

type RawLog struct {
	gorm.Model
	ComponentID  uint    `json:"component_id"`
	MetricTypeID uint    `json:"metric_type_id"`
	Value        float64 `json:"value"`
	Timestamp    string  `json:"timestamp"`

	Component  Component  `gorm:"foreignKey:ComponentID"`
	MetricType MetricType `gorm:"foreignKey:MetricTypeID"`
}
