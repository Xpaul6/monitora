package limitutils

import (
	"log"
	"time"

	. "github.com/XPaul6/monitora/models"
	"gorm.io/gorm"
)

func DetectLimitCross(logPair IdInfoPair, metricTypeMap map[uint]MetricType, db *gorm.DB) {
	var limits []Limit
	result := db.Where("component_id IN (?)",
		db.Table("components").Where("server_id = ?", logPair.Id).Select("id"),
	).Find(&limits)
	if result.Error != nil {
		log.Printf("Error on fetching limits: %v", result.Error)
		return
	}

	if len(limits) == 0 {
		return
	}

	var components []Component
	db.Where("server_id = ?", logPair.Id).Find(&components)
	componentMap := make(map[uint]Component)
	for _, c := range components {
		componentMap[c.ID] = c
	}

	var notifications []Notification
	for _, limit := range limits {
		comp, ok := componentMap[limit.ComponentID]
		if !ok {
			continue
		}

		metricType, ok := metricTypeMap[limit.MetricTypeID]
		if !ok {
			continue
		}

		value := getCurrentValue(logPair.Info, comp.Type, comp.Address, metricType.Name)
		if value == 0 {
			continue
		}

		if value > limit.ThresholdValue {
			notifications = append(notifications, Notification{
				LimitID:   limit.ID,
				RealValue: value,
				Timestamp: time.Now(),
			})
		}
	}
	if err := db.Create(&notifications).Error; err != nil {
		log.Printf("Error creating notifications: %v", err)
	}
}

func getCurrentValue(info SysInfo, compType, address, metricName string) float64 {
	switch compType {
	case "cpu":
		switch metricName {
		case "cpu_load":
			return info.CPU.LoadPercentage
		case "cpu_temperature":
			return info.CPU.Temperature
		}
	case "mem":
		switch metricName {
		case "mem_load":
			return info.Mem.LoadPercentage
		case "mem_total":
			return float64(info.Mem.Total)
		case "mem_used":
			return float64(info.Mem.Used)
		}
	case "disk":
		for _, d := range info.Disks {
			if d.MountPoint == address {
				switch metricName {
				case "disk_total":
					return float64(d.Total)
				case "disk_used":
					return float64(d.Used)
				}
			}
		}
	case "net":
		for _, n := range info.Net {
			if n.Name == address {
				switch metricName {
				case "net_rbps":
					return n.RBpS
				case "net_sbps":
					return n.SBpS
				}
			}
		}
	}
	return 0
}
