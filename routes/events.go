package routes

import (
	"net/http"
	"strconv"

	"example.com/event/handlers"
	"example.com/event/models"
	"github.com/gin-gonic/gin"
)

func getEvents(context *gin.Context) {
	events, err := handlers.GetAllEvents()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not retrieve events",
			"error":   err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "List of events",
		"data":    events,
	})
}

func getEventById(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("eventId"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse event id",
			"error":   err.Error(),
		})
		return
	}

	event, err := handlers.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"message": "Event not found",
			"error":   err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "Successfully get event details",
		"data":    event,
	})
}

func createEvent(context *gin.Context) {
	var newEvent models.Event
	err := context.ShouldBindJSON(&newEvent)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse request",
			"error":   err.Error(),
		})
		return
	}

	// Retrieve userId from request context (set earlier by auth middleware)
	userId := context.GetInt64("userId")

	newEvent.UserID = userId

	err = handlers.CreateEvent(&newEvent)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not save event",
			"error":   err.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Event created successfully",
		"data":    newEvent,
	})
}

func updateEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("eventId"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid event ID",
			"error":   err.Error(),
		})
		return
	}

	var updatedEvent models.Event
	err = context.ShouldBindJSON(&updatedEvent)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse request",
			"error":   err.Error(),
		})
		return
	}

	// Check if the event exists
	event, err := handlers.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"message": "Event not found",
			"error":   err.Error(),
		})
		return
	}

	// Event only can be updated by event creator
	userId := context.GetInt64("userId")
	if event.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Not authorized to update event",
			"error":   "not authorized to update event",
		})
		return
	}

	updatedEvent.ID = eventId
	updatedEvent.UserID = userId

	err = handlers.UpdateEvent(&updatedEvent)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not update event",
			"error":   err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Event updated successfully",
		"data":    updatedEvent,
	})
}

func deleteEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("eventId"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid event ID",
			"error":   err.Error(),
		})
		return
	}

	// Check if the event exists
	event, err := handlers.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"message": "Event not found",
			"error":   err.Error(),
		})
		return
	}

	userId := context.GetInt64("userId")

	// Event only can be deleted by event creator
	if event.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Not authorized to delete event",
			"error":   "not authorized to delete event",
		})
		return
	}

	err = handlers.DeleteEvent(event)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not delete event",
			"error":   err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Event deleted successfully",
	})
}
