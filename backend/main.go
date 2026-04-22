package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/XPaul6/monitora/controllers"
	authutils "github.com/XPaul6/monitora/utils/auth"
	dbutil "github.com/XPaul6/monitora/utils/database"
	fetchutil "github.com/XPaul6/monitora/utils/fetch"
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

	// Information gathering
	go fetchutil.RunFetchUtil(db)

	// Router setup
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	auth := router.Group("/auth")
	{
		auth.POST("/register", controllers.Register(db))
		auth.POST("/login", controllers.Login(db))
		auth.GET("/check", authutils.WithAuth(db), func(c *gin.Context) { c.JSON(200, gin.H{"status": "authorized"}) })
	}

	user := router.Group("/user")
	{
		// Servers
		user.GET("/servers", authutils.WithAuth(db), controllers.GetAllServers(db))
		user.GET("/components", authutils.WithAuth(db), controllers.GetServerComponents(db))
		user.POST("/add-server", authutils.WithAuth(db), controllers.AddServer(db))
		user.DELETE("/delete-server", authutils.WithAuth(db), controllers.DeleteServer(db))

		// Limits
		user.GET("/limits", authutils.WithAuth(db), controllers.GetLimits(db))
		user.POST("/set-limit", authutils.WithAuth(db), controllers.SetLimit(db))
		user.DELETE("/delete-limit", authutils.WithAuth(db), controllers.DeleteLimit(db))

		// Common getters
		user.GET("/notifications", authutils.WithAuth(db), controllers.GetNotifications(db))
	}

	stats := router.Group("/stats")
	{
		stats.GET("/by-period", authutils.WithAuth(db), controllers.GetStatsByPeriod(db))
	}

	service := router.Group("/service")
	{
		service.GET("/metric-types", controllers.GetMetricTypes(db))
	}

	router.Run(":8080")
}
