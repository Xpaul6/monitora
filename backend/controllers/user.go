package controllers

import (
	"net/http"

	. "github.com/XPaul6/monitora/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllServers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userVal, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied"})
			return
		}
		user := userVal.(User)

		var serverList []Server
		result := db.Find(&serverList, "user_id = ?", user.ID)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
			return
		}

		resBody := GetAllServersResponse{
			Count: result.RowsAffected,
			Servers: serverList,
		}

		c.JSON(http.StatusOK, resBody)
	}
}

func AddServer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody AddServerRequest
		if err := c.Bind(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		userVal, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied"})
			return
		}
		user := userVal.(User)

		var num int64
		db.Model(&Server{}).Where("ip = ? and user_id = ?", reqBody.IP, user.ID).Count(&num)
		if num > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "server with this ip already exists"})
			return
		}

		newServer := Server{
			Name: reqBody.Name,
			IP: reqBody.IP,
			UserID: user.ID,
		}

		result := db.Create(&newServer)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
			return
		}
		c.JSON(http.StatusOK, gin.H{"server_id": newServer.ID})
	}
}

func DeleteServer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody DeleteServerRequest
		if err := c.Bind(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		userVal, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied"})
			return
		}
		user := userVal.(User)

		var server Server
		result := db.First(&server, reqBody.ID)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
			return
		}

		if server.UserID != user.ID {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied"})
			return
		}

		result = db.Delete(&server)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}
