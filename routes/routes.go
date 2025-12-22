package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	// GET All Events
	server.GET("/events", getEvents)

	// GET Event by ID
	server.GET("/event/:eventId", getEventById)

	// POST Create Event
	server.POST("/event", createEvent)

	// PUT Update Event
	server.PUT("/event/:eventId", updateEvent)

	// DELETE Event
	server.DELETE("/event/:eventId", deleteEvent)

	// Sign Up User
	server.POST("/user/signup", signUp)

	// Login User
	server.POST("/user/login", login)
}
