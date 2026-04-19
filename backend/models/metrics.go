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
// (cpu_load,%,"Cpu time occupancy percentage")
// (cpu_temperature,celsius,"CPU package temperature")
// (mem_load,%,"RAM load percentage")
// (mem_total,bytes,"Total RAM avaliable")
// (mem_used,bytes,"Used RAM")
// (disk_mountpoint,string,"Disk device mount point")
// (disk_total,bytes,"Total disk device capacity")
// (disk_used,bytes,"Amount of used disk space")
// (net_name,string,"Network device name")
// (net_rbps,bps,"Recieved bytes per second")
// (net_sbps,bps,"Sent bytes per second")

 var DefaultTypes []MetricType

 func init() {
	DefaultTypes = []MetricType{
		{
			Name: "cpu_load",
			Unit: "%",
			Description: "Cpu time occupancy percentage",
		},
		{
			Name: "cpu_temperature",
			Unit: "celsius",
			Description: "CPU package temperature",
		},
		{
			Name: "mem_load",
			Unit: "%",
			Description: "RAM load percentage",
		},
		{
			Name: "mem_total",
			Unit: "bytes",
			Description: "Total RAM avaliable",
		},
		{
			Name: "mem_used",
			Unit: "bytes",
			Description: "Used RAM",
		},
		{
			Name: "disk_mountpoint",
			Unit: "string",
			Description: "Disk device mount point",
		},
		{
			Name: "disk_total",
			Unit: "bytes",
			Description: "Total disk device capacity",
		},
		{
			Name: "disk_used",
			Unit: "bytes",
			Description: "Amount of used disk space",
		},
		{
			Name: "net_name",
			Unit: "string",
			Description: "Network device name",
		},
		{
			Name: "net_rbps",
			Unit: "bps",
			Description: "Recieved bytes per second",
		},
		{
			Name: "net_sbps",
			Unit: "bps",
			Description: "Sent bytes per second",
		},
	}
 }
