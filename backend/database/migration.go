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

	res := db.Raw(`
		SELECT create_hypertable(
			'raw_metrics',
			'timestamp',
			chunk_time_interval => INTERVAL '1 day',
			if_not_exists => true
		);
	`)
	if res.Error != nil {
		return err
	}

	return nil
}
