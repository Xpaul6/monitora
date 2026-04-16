package models

import (
	"time"
)

type RawLog struct {
	ComponentID  uint      `json:"component_id"`
	MetricTypeID uint      `json:"metric_type_id"`
	Value        float64   `json:"value"`
	Timestamp    time.Time `json:"timestamp"`

	Component  Component  `gorm:"foreignKey:ComponentID"`
	MetricType MetricType `gorm:"foreignKey:MetricTypeID"`
}
