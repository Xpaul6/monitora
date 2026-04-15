package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/XPaul6/monitora/database"
	. "github.com/XPaul6/monitora/models"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Cannot load .env file, using default vars")
	}

	var dbconf DBConfig = DBConfig{
		Host: os.Getenv("DB_HOST"),
		User: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName: os.Getenv("DB_NAME"),
		Port: os.Getenv("DB_PORT"),
	}

	db, err := dbutil.CreateDBConnection(dbconf)
	if err != nil {
		log.Fatalln("Cannot connect to the database")
	}

	fmt.Println(db.Migrator().HasTable("Users"))

	router := gin.Default()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})

	router.Run("localhost:8080")
}
