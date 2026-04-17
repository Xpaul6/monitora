package models

type DBConfig struct {
	Host     string `default:"postgres"`
	User     string `default:"postgres"`
	Password string `default:"postgres"`
	DBName   string `default:"monitora"`
	Port     string `default:"5432"`
}
