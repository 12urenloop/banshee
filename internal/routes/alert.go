package routes

import (
	"12ul/banshee/internal/alerts"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	alerts := rg.Group("/alerts")
	{
		alerts.GET("/", getAlerts)
		alerts.PUT("/:alertid", dismissAlert)
	}
}

func getAlerts(c *gin.Context) {
	dismissedIds := alerts.GetDismissedAlerts()
	unDismissedAlerts := alerts.GetUndismissedAlerts()
	c.JSON(200, gin.H{
		"alerts":          unDismissedAlerts,
		"dismissedAlerts": dismissedIds,
	})
}

func dismissAlert(c *gin.Context) {
	id := c.Param("alertid")
	alerts.DismissAlert(id)
	c.JSON(200, gin.H{})
}
