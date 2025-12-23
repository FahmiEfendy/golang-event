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

	// Register Event
	server.POST("/event/:id/register", registerEvent)

	// Unregister Event
	server.DELETE("/event/:id/unregister", unregisterEvent)

	// Sign Up User
	server.POST("/user/signup", signUp)

	// Login User
	server.POST("/user/login", login)

	// TODO: GET Registration List

	// Other way to register protected routes
	// authenticated := server.Group("/")
	// authenticated.Use(middlewares.Authenticate)
	// authenticated.POST("/event", createEvent)
	// authenticated.PUT("/events/:id", updateEvent)
	// authenticated.DELETE("/events/:id", deleteEvent)
	// authenticated.POST("/event/:id/register", registerEvent)
	// authenticated.DELETE("/event/:id/unregister", unregisterEvent)
}
