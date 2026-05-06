package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	. "github.com/XPaul6/monitora/models"
)

func GetLimits(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody GetLimitsRequest
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

		var limits []Limit
		result := db.Table("limits").
			Joins("JOIN components ON components.id = limits.component_id").
			Joins("JOIN servers ON servers.id = components.server_id").
			Where("servers.user_id = ? AND servers.id = ?", user.ID, reqBody.ID).
			Find(&limits)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
			return
		}

		c.JSON(http.StatusOK, limits)
	}
}

func SetLimit(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody SetLimitRequest
		if err := c.Bind(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		// Ownership checks
		userVal, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied"})
			return
		}
		user := userVal.(User)
		var component Component
		if err := db.First(&component, reqBody.ComponentID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
			return
		}
		var server Server
		if err := db.First(&server, component.ServerID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
			return
		}
		if server.UserID != user.ID {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied"})
			return
		}

		newLimit := Limit{
			ComponentID:    reqBody.ComponentID,
			MetricTypeID:   reqBody.MetricTypeID,
			ThresholdValue: reqBody.ThresholdValue,
		}

		result := db.Create(&newLimit)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": newLimit.ID})
	}
}

func DeleteLimit(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody DeleteLimitRequest
		if err := c.Bind(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		// Ownership checks
		userVal, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied"})
			return
		}
		user := userVal.(User)
		var limit Limit
		if err := db.First(&limit, reqBody.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
			return
		}
		var component Component
		if err := db.First(&component, limit.ComponentID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
			return
		}
		var server Server
		if err := db.First(&server, component.ServerID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
			return
		}
		if server.UserID != user.ID {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied"})
			return
		}

		result := db.Delete(&Limit{}, reqBody.ID)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}
