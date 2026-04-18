package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/XPaul6/monitora/controllers"
	authutils "github.com/XPaul6/monitora/utils/auth"
	dbutil "github.com/XPaul6/monitora/utils/database"
)

func init() {
	// gin.SetMode(gin.ReleaseMode)
	err := godotenv.Load()
	if err != nil {
		log.Println("Cannot load .env file, using default vars")
	}
}

func main() {
	// Database connection
	db, err := dbutil.CreateDBConnection()
	if err != nil {
		log.Fatalln("Cannot connect to the database")
	}
	err = dbutil.AutoMigrate(db)
	if err != nil {
		log.Fatalln("Cannot migrate to database")
	}

	// Router setup
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	router.GET("/users", controllers.GetAllUsers(db))

	auth := router.Group("/auth")
	{
		auth.POST("/register", controllers.Register(db))
		auth.POST("/login", controllers.Login(db))
		auth.GET("/check", authutils.WithAuth(db), func(c *gin.Context) { c.JSON(200, gin.H{"status": "authorized"}) })
	}

	router.Run(":8080")
}
