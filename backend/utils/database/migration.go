package dbutil

import (
	. "github.com/XPaul6/monitora/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&User{},
		&Server{},
		&MetricType{},
		&Component{},
		&Limit{},
		&Notification{},
		&RawLog{},
	)
	if err != nil {
		return err
	}

	create_ht_sql := `
		SELECT create_hypertable(
			'raw_logs',
			'timestamp',
			chunk_time_interval => INTERVAL '1 day',
			if_not_exists => true
		);
	`
	res := db.Exec(create_ht_sql)
	if res.Error != nil {
		return res.Error
	}

	var count int64
	db.Model(&MetricType{}).Count(&count)
	if count == 0 {
		db.Create(&DefaultTypes)
	}

	return nil
}
