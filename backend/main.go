package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/XPaul6/monitora/database"
	// . "github.com/XPaul6/monitora/models"
)

func init() {

}
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Cannot load .env file, using default vars")
	}

	db, err := dbutil.CreateDBConnection()
	if err != nil {
		log.Fatalln("Cannot connect to the database")
	}
	err = dbutil.AutoMigrate(db)
	if err != nil {
		log.Fatalln("Cannot migrate to database")
	}

	fmt.Println(db.Migrator().HasTable("Users"))

	router := gin.Default()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})

	router.Run("localhost:8080")
}
