package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Setup engine (configure HTTP server)
	server := gin.Default()

	// GET "/"
	server.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "Successfully connected to the server",
		})
	})

	// Start server on localhost:8080
	server.Run(":8080")
}
