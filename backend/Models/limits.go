package models

import (
	"gorm.io/gorm"
)

type Limit struct {
	gorm.Model
	ComponentID    uint    `json:"component_id"`
	MetricTypeID   uint    `json:"metric_type_id"`
	ThresholdValue float64 `json:"threshhold_value"`

	Server Server `gorm:"foreignKey:ServerID"`
	Component Component `gorm:"foreignKey:ComponentID"`
	MetricType MetricType `gorm:"foreignKey:MetricTypeID"`
	Notifications []Notification `gorm:"foreignKey:LimitID"`
}
