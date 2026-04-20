package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	. "github.com/XPaul6/monitora/models"
)

func GetStatsByPeriod(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody GetStatsByPeriodRequest
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userVal, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied"})
			return
		}
		user := userVal.(User)

		var server Server
		result := db.First(&server, reqBody.ServerID)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
			return
		}

		if server.UserID != user.ID {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied"})
			return
		}

		var logs []RawLog
		result = db.Where("timestamp >= ? and timestamp <= ?", reqBody.PeriodBegin, reqBody.PeriodEnd).Find(&logs)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		}

		// Mertic types caching
		var metricTypes []MetricType
		db.Find(&metricTypes)
		typeMap := make(map[uint]MetricType)
		for _, mt := range metricTypes {
			typeMap[mt.ID] = mt
		}

		// Components cahcing
		var serverComponents []Component
		db.Where("server_id = ?", reqBody.ServerID).Find(&serverComponents)
		componentMap := make(map[uint]Component)
		for _, c := range serverComponents {
			componentMap[c.ID] = c
		}

		// Forming response
		var returnLogs []GetStatsByPeriodResponse
		for _, log := range logs {
			returnLogs = append(returnLogs, GetStatsByPeriodResponse{
				Component:  componentMap[log.ComponentID],
				MetricType: typeMap[log.MetricTypeID],
				Value:      log.Value,
				TimeStamp:  log.Timestamp,
			})
		}

		c.JSON(http.StatusOK, returnLogs)
	}
}
