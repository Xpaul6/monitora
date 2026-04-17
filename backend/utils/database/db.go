package dbutil

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CreateDBConnection() (*gorm.DB, error) {
	var dsn string = os.Getenv("DB_CONFIG")
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
