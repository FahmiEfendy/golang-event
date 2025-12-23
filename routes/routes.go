package routes

import (
	"example.com/event/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	// GET All Events
	server.GET("/events", getEvents)

	// GET Event by ID
	server.GET("/event/:eventId", getEventById)

	// POST Create Event
	server.POST("/event", middlewares.Authenticate, createEvent)

	// PUT Update Event
	server.PUT("/event/:eventId", middlewares.Authenticate, updateEvent)

	// DELETE Event
	server.DELETE("/event/:eventId", middlewares.Authenticate, deleteEvent)

	// Sign Up User
	server.POST("/user/signup", signUp)

	// Login User
	server.POST("/user/login", login)

	// Other way to register protected routes
	// authenticated := server.Group("/")
	// authenticated.Use(middlewares.Authenticate)
	// authenticated.POST("/event", createEvent)
	// authenticated.PUT("/events/:id", updateEvent)
	// authenticated.DELETE("/events/:id", deleteEvent)
}
