package main

import (
	"12ul/banshee/internal/alerts"
	"12ul/banshee/internal/config"
	"12ul/banshee/internal/routes"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Create a new gin router
	r := gin.Default()

	// All files besides the index page are accessible via /public endpoint
	r.StaticFS("/public", http.Dir("public"))
	r.StaticFile("/", "public/index.html")

	alerts.StartFetchInterval()

	r.Use(
		cors.New(cors.Config{
			AllowAllOrigins: true,
		}),
		gin.Recovery(),
	)

	// Register routes
	rg := r.Group("/api/v1")
	routes.RegisterRoutes(rg)

	r.Run(":" + config.GetEnv("PORT"))
}
