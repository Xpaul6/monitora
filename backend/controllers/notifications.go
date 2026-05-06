package controllers

import (
	"net/http"

	. "github.com/XPaul6/monitora/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetNotifications(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody GetNotificationsRequest
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
		var notifications []Notification
		result := db.Table("notifications").
			Joins("JOIN limits ON limits.id = notifications.limit_id").
			Joins("JOIN components ON components.id = limits.component_id").
			Joins("JOIN servers ON servers.id = components.server_id").
			Where("servers.user_id = ? AND servers.id = ?", user.ID, reqBody.ID).
			Find(&notifications)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
			return
		}

		c.JSON(http.StatusOK, notifications)
	}
}
