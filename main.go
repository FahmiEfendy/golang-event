package main

import (
	"net/http"

	"example.com/event/db"
	"example.com/event/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	db.InitDB()

	// Setup engine (configure HTTP server)
	server := gin.Default()

	// GET "/"
	server.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "Successfully connected to the server",
		})
	})

	routes.RegisterRoutes(server)

	// Start server on localhost:8080
	server.Run(":8080")
}
