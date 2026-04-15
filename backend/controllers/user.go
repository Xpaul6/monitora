package controllers

import (
	. "github.com/XPaul6/monitora/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		users := []User{}
		result := db.Find(&users)
		if result.Error != nil {
			c.Status(500)
			return
		}

		c.IndentedJSON(200, users)
	}
}
