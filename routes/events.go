package routes

import (
	"net/http"
	"strconv"

	"example.com/event/models"
	"github.com/gin-gonic/gin"
)

func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()
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
			"message": "Invalid event ID",
			"error":   err.Error(),
		})
		return
	}

	event, err := models.GetEventByID(eventId)
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

	// Dummy assignment of ID and UserID
	newEvent.ID = 1
	newEvent.UserID = 1

	err = newEvent.Save()
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

	// Check if the event exists
	_, err = models.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"message": "Event not found",
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

	updatedEvent.ID = eventId

	err = updatedEvent.Update()
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
	event, err := models.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"message": "Event not found",
			"error":   err.Error(),
		})
		return
	}

	err = event.Delete()
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
