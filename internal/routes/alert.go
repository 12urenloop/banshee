package routes

import (
	"12ul/banshee/internal/alerts"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	alerts := rg.Group("/alerts")
	{
		alerts.GET("/", getAlerts)
		alerts.POST("/dismiss", dismissAlert)
	}
}

func getAlerts(c *gin.Context) {
	alerts := alerts.GetUndismissedAlerts()
	c.JSON(200, gin.H{
		"alerts": alerts,
	})
}

func dismissAlert(c *gin.Context) {
	var reqBody struct {
		Id string `json:"id"`
	}
	c.BindJSON(&reqBody)
	alerts.DismissAlert(reqBody.Id)
	c.JSON(200, gin.H{})
}
