package middlewares

import (
	"net/http"

	"example.com/event/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")

	// Check if token not in headers
	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Not authorized",
		})
		return
	}

	// Verify token
	userId, err := utils.VerifyToken(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Could not verify token",
			"error":   err.Error(),
		})
		return
	}

	// Store userId in request context so it can be accessed by next handlers
	context.Set("userId", userId)

	// Continue to next handlers
	context.Next()
}
