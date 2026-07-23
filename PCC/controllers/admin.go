package controllers

import (
	"net/http"

	"main/wa"

	"github.com/gin-gonic/gin"
)

// ReloadSchedule handles admin request to reload schedule cache
func ReloadSchedule(c *gin.Context) {
	wa.ReloadScheduleCache()
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "Schedule cache reloaded"})
}
