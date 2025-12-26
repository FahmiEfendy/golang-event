package routes

import (
	"bytes"
	"encoding/json"
	"testing"

	"example.com/event/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// USER ROUTES
const LOGIN_PATH = "/login"
const SIGNUP_PATH = "/signup"

// EVENT ROUTES
const GET_EVENTS_PATH = "/events"
const GET_EVENTS_BY_ID_PATH = "/event/:eventId"
const CREATE_EVENT_PATH = "/event"
const UPDATE_EVENT_PATH = "/event/:eventId"
const DELETE_EVENT_PATH = "/event/:eventId"

// convert []byte (payload) to JSON
func toJSON(t *testing.T, v any) *bytes.Buffer {
	b, err := json.Marshal(v)
	assert.NoError(t, err)
	return bytes.NewBuffer(b)
}

// creates and returns a Gin HTTP router that is configured specifically for tests
func setupRouter() *gin.Engine {
	// indicate doing test mode
	gin.SetMode(gin.TestMode)

	r := gin.Default()

	// define user routes
	r.POST(SIGNUP_PATH, signUp)
	r.POST(LOGIN_PATH, login)

	// define event routes
	r.GET(GET_EVENTS_PATH, getEvents)
	r.GET(GET_EVENTS_BY_ID_PATH, getEventById)
	r.POST(CREATE_EVENT_PATH, createEvent)
	r.PUT(UPDATE_EVENT_PATH, middlewares.Authenticate, updateEvent)
	r.DELETE(UPDATE_EVENT_PATH, middlewares.Authenticate, deleteEvent)

	return r
}
