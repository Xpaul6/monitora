package controllers

import (
	"net/http"

	. "github.com/XPaul6/monitora/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetMetricTypes(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var metricTypes []MetricType
		result := db.Find(&metricTypes)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
			return
		}

		c.JSON(http.StatusOK, metricTypes)
	}
}
