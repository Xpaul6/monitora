package dbutil

import (
	"fmt"

	. "github.com/XPaul6/monitora/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CreateDBConnection(cfg DBConfig) (*gorm.DB, error) {
	var dsn string = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port,
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
