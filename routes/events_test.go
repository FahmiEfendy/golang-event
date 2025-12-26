package routes

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"example.com/event/handlers"
	"example.com/event/models"
	"example.com/event/utils"
	"github.com/stretchr/testify/assert"
)

var createEventPayload = map[string]string{
	"name":        "Go Workshop Jakarta",
	"description": "A beginner-friendly workshop covering Go fundamentals and best practices.",
	"location":    "Jakarta",
	"dateTime":    "2025-12-16T09:00:00+07:00",
}

var updateEventPayload = map[string]string{
	"name":        "Go Workshop Bandung",
	"description": "A beginner-friendly workshop covering Go fundamentals.",
	"location":    "Bandung",
	"dateTime":    "2025-12-26T09:00:00+07:00",
}

func TestGetEvents_Success(t *testing.T) {
	router := setupRouter()

	// keep original handler to restore global state after test
	originalGetAllEvents := handlers.GetAllEvents

	// restore original handler after test to avoid side effects
	defer func() {
		handlers.GetAllEvents = originalGetAllEvents
	}()

	// mock handlers.GetAllEvents
	handlers.GetAllEvents = func() ([]models.Event, error) {
		return []models.Event{
				{
					ID:          1,
					Name:        "Go Workshop Jakarta",
					Description: "A beginner-friendly workshop covering Go fundamentals and best practices.",
					Location:    "Jakarta",
					DateTime:    time.Date(2025, 12, 16, 9, 0, 0, 0, time.FixedZone("WIB", 7*3600)),
					UserID:      1,
				},
			},
			nil
	}

	// simulate hit API
	req, _ := http.NewRequest(
		http.MethodGet,
		GET_EVENTS_PATH,
		http.NoBody,
	)

	// using w to simulate expected response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// expected response statusCode
	assert.Equal(t, http.StatusOK, w.Code)

	// convert JSON to []byte
	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// expected response data
	assert.Equal(t, "List of events", response["message"])
	assert.NotEmpty(t, response["data"])

}

func TestGetEvents_ErrorGetAllEvents(t *testing.T) {
	router := setupRouter()

	originalGetAllEvents := handlers.GetAllEvents
	defer func() {
		handlers.GetAllEvents = originalGetAllEvents
	}()
	handlers.GetAllEvents = func() ([]models.Event, error) {
		return []models.Event{{}}, errors.New("simulate error get all events")
	}

	req, _ := http.NewRequest(
		http.MethodGet,
		GET_EVENTS_PATH,
		http.NoBody,
	)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Could not retrieve events", response["message"])
	assert.NotEmpty(t, response["error"])

}

func TestGetEventById_Success(t *testing.T) {
	router := setupRouter()

	// keep original handler to restore global state after test
	originalGetEventByID := handlers.GetEventByID

	// restore original handler after test to avoid side effects
	defer func() {
		handlers.GetEventByID = originalGetEventByID
	}()

	// mock handlers.GetEventByID
	handlers.GetEventByID = func(eventId int64) (*models.Event, error) {
		return &models.Event{
				ID:          1,
				Name:        "Go Workshop Jakarta",
				Description: "A beginner-friendly workshop covering Go fundamentals and best practices.",
				Location:    "Jakarta",
				DateTime:    time.Date(2025, 12, 16, 9, 0, 0, 0, time.FixedZone("WIB", 7*3600)),
				UserID:      1,
			},
			nil
	}

	// simulate hit API
	req, _ := http.NewRequest(
		http.MethodGet,
		strings.Replace(GET_EVENTS_BY_ID_PATH, ":eventId", "1", 1),
		http.NoBody,
	)

	// using w to simulate expected response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// expected response statusCode
	assert.Equal(t, http.StatusOK, w.Code)

	// convert JSON to []byte
	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// expected response data
	assert.Equal(t, "Successfully get event details", response["message"])
	assert.NotEmpty(t, response["data"])

}

func TestGetEventById_ErrorParseEventId(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest(
		http.MethodGet,
		strings.Replace(GET_EVENTS_BY_ID_PATH, ":eventId", "invalidEventId", 1),
		http.NoBody,
	)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Could not parse event id", response["message"])
	assert.NotEmpty(t, response["error"])

}

func TestGetEventById_ErrorGetEventById(t *testing.T) {
	router := setupRouter()

	originalGetEventByID := handlers.GetEventByID
	defer func() {
		handlers.GetEventByID = originalGetEventByID
	}()
	handlers.GetEventByID = func(eventId int64) (*models.Event, error) {
		return &models.Event{},
			errors.New("simulate error get event by id")
	}

	req, _ := http.NewRequest(
		http.MethodGet,
		strings.Replace(GET_EVENTS_BY_ID_PATH, ":eventId", "67", 1),
		http.NoBody,
	)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Event not found", response["message"])
	assert.NotEmpty(t, response["error"])
}

func TestCreateEvent_Success(t *testing.T) {
	router := setupRouter()

	body := toJSON(t, createEventPayload)

	// keep original handler to restore global state after test
	originalCreateEvent := handlers.CreateEvent

	// restore original handler after test to avoid side effects
	defer func() {
		handlers.CreateEvent = originalCreateEvent
	}()

	// mock handlers.CreateEvent
	handlers.CreateEvent = func(event *models.Event) error {
		event.ID = 1 // simulate DB insert
		return nil
	}

	// simulate hit API
	req, _ := http.NewRequest(
		http.MethodPost,
		CREATE_EVENT_PATH,
		body,
	)

	// using w to simulate expected response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// expected response statusCode
	assert.Equal(t, http.StatusCreated, w.Code)

	// convert JSON to []byte
	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// map response data
	data := response["data"].(map[string]any)

	// expected response data
	assert.Equal(t, "Event created successfully", response["message"])
	assert.Equal(t, createEventPayload["name"], data["Name"])
	assert.Equal(t, createEventPayload["description"], data["Description"])
	assert.Equal(t, createEventPayload["location"], data["Location"])
	assert.Equal(t, createEventPayload["dateTime"], data["DateTime"])
	assert.Equal(t, float64(1), data["ID"]) // JSON numbers → float64
}

func TestCreateEvent_ErrorShouldBindJSON(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest(
		http.MethodPost,
		CREATE_EVENT_PATH,
		bytes.NewBufferString("simulate invalid JSON"),
	)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Could not parse request", response["message"])
	assert.NotEmpty(t, response["error"].(string))

}

func TestCreateEvent_ErrorCreateEventHandler(t *testing.T) {
	router := setupRouter()

	body := toJSON(t, createEventPayload)

	originalCreateEvent := handlers.CreateEvent
	defer func() {
		handlers.CreateEvent = originalCreateEvent
	}()
	handlers.CreateEvent = func(event *models.Event) error {
		return errors.New("simulate error create event handler")
	}

	req, _ := http.NewRequest(
		http.MethodPost,
		CREATE_EVENT_PATH,
		body,
	)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Could not save event", response["message"])
	assert.NotEmpty(t, response["error"].(string))
}

func TestUpdateEvent_Success(t *testing.T) {
	router := setupRouter()

	body := toJSON(t, updateEventPayload)

	// mock utils.VerifyToken
	originalVerifyToken := utils.VerifyToken
	defer func() {
		utils.VerifyToken = originalVerifyToken
	}()
	utils.VerifyToken = func(token string) (int64, error) {
		return 1, nil
	}

	// mock handlers.GetEventByID
	originalGetEventByID := handlers.GetEventByID
	defer func() {
		handlers.GetEventByID = originalGetEventByID
	}()
	handlers.GetEventByID = func(eventId int64) (*models.Event, error) {
		return &models.Event{
				ID:          1,
				Name:        "Go Workshop Jakarta",
				Description: "A beginner-friendly workshop covering Go fundamentals and best practices.",
				Location:    "Jakarta",
				DateTime:    time.Date(2025, 12, 16, 9, 0, 0, 0, time.FixedZone("WIB", 7*3600)),
				UserID:      1,
			},
			nil
	}

	// mock handlers.UpdateEvent
	originalUpdateEvent := handlers.UpdateEvent
	defer func() {
		handlers.UpdateEvent = originalUpdateEvent
	}()
	handlers.UpdateEvent = func(event *models.Event) error {
		event.ID = 1
		return nil
	}

	// simulate hit API
	req, _ := http.NewRequest(
		http.MethodPut,
		strings.Replace(UPDATE_EVENT_PATH, ":eventId", "1", 1),
		body,
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "sample-token-userid-1")

	// using w to simulate expected response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// expected response statusCode
	assert.Equal(t, http.StatusOK, w.Code)

	// convert JSON to []byte
	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// map response data
	data := response["data"].(map[string]any)

	// expected response data
	assert.Equal(t, "Event updated successfully", response["message"])
	assert.Equal(t, float64(1), data["ID"])
	assert.Equal(t, updateEventPayload["name"], data["Name"])
	assert.Equal(t, updateEventPayload["description"], data["Description"])
	assert.Equal(t, updateEventPayload["location"], data["Location"])
	assert.Equal(t, updateEventPayload["dateTime"], data["DateTime"])
	assert.Equal(t, float64(1), data["UserID"])
}

func TestUpdateEvent_ErrorParseEventId(t *testing.T) {
	router := setupRouter()

	body := toJSON(t, updateEventPayload)

	originalVerifyToken := utils.VerifyToken
	defer func() {
		utils.VerifyToken = originalVerifyToken
	}()
	utils.VerifyToken = func(token string) (int64, error) {
		return 1, nil
	}

	req, _ := http.NewRequest(
		http.MethodPut,
		strings.Replace(UPDATE_EVENT_PATH, ":eventId", "invalidEventId", 1),
		body,
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "sample-token-userid-1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Invalid event ID", response["message"])
	assert.NotEmpty(t, response["error"].(string))
}

func TestUpdateEvent_ErrorGetEventById(t *testing.T) {
	router := setupRouter()

	body := toJSON(t, updateEventPayload)

	originalVerifyToken := utils.VerifyToken
	defer func() {
		utils.VerifyToken = originalVerifyToken
	}()
	utils.VerifyToken = func(token string) (int64, error) {
		return 1, nil
	}

	originalGetEventByID := handlers.GetEventByID
	defer func() {
		handlers.GetEventByID = originalGetEventByID
	}()
	handlers.GetEventByID = func(eventId int64) (*models.Event, error) {
		return &models.Event{},
			errors.New("simulate error get event by id")
	}

	req, _ := http.NewRequest(
		http.MethodPut,
		strings.Replace(UPDATE_EVENT_PATH, ":eventId", "1", 1),
		body,
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "sample-token-userid-1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Event not found", response["message"])
	assert.NotEmpty(t, response["error"].(string))
}

func TestUpdateEvent_ErrorUnauthorize(t *testing.T) {
	router := setupRouter()

	body := toJSON(t, updateEventPayload)

	originalVerifyToken := utils.VerifyToken
	defer func() {
		utils.VerifyToken = originalVerifyToken
	}()
	utils.VerifyToken = func(token string) (int64, error) {
		return 67, nil
	}

	originalGetEventByID := handlers.GetEventByID
	defer func() {
		handlers.GetEventByID = originalGetEventByID
	}()
	handlers.GetEventByID = func(eventId int64) (*models.Event, error) {
		return &models.Event{
				ID:          1,
				Name:        "Go Workshop Jakarta",
				Description: "A beginner-friendly workshop covering Go fundamentals and best practices.",
				Location:    "Jakarta",
				DateTime:    time.Date(2025, 12, 16, 9, 0, 0, 0, time.FixedZone("WIB", 7*3600)),
				UserID:      1,
			},
			nil
	}

	req, _ := http.NewRequest(
		http.MethodPut,
		strings.Replace(UPDATE_EVENT_PATH, ":eventId", "1", 1),
		body,
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "sample-token-userid-67")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Not authorized to update event", response["message"])
	assert.Equal(t, "not authorized to update event", response["error"])
}

func TestUpdateEvent_ErrorShouldBindJSON(t *testing.T) {
	router := setupRouter()

	originalVerifyToken := utils.VerifyToken
	defer func() {
		utils.VerifyToken = originalVerifyToken
	}()
	utils.VerifyToken = func(token string) (int64, error) {
		return 1, nil
	}

	req, _ := http.NewRequest(
		http.MethodPut,
		strings.Replace(UPDATE_EVENT_PATH, ":eventId", "1", 1),
		bytes.NewBufferString("simulate invalid JSON"),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "sample-token-userid-1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Could not parse request", response["message"])
	assert.NotEmpty(t, response["error"].(string))
}

func TestUpdateEvent_ErrorUpdateEventHandler(t *testing.T) {
	router := setupRouter()

	body := toJSON(t, updateEventPayload)

	originalVerifyToken := utils.VerifyToken
	defer func() {
		utils.VerifyToken = originalVerifyToken
	}()
	utils.VerifyToken = func(token string) (int64, error) {
		return 1, nil
	}

	originalGetEventByID := handlers.GetEventByID
	defer func() {
		handlers.GetEventByID = originalGetEventByID
	}()
	handlers.GetEventByID = func(eventId int64) (*models.Event, error) {
		return &models.Event{
				ID:          1,
				Name:        "Go Workshop Jakarta",
				Description: "A beginner-friendly workshop covering Go fundamentals and best practices.",
				Location:    "Jakarta",
				DateTime:    time.Date(2025, 12, 16, 9, 0, 0, 0, time.FixedZone("WIB", 7*3600)),
				UserID:      1,
			},
			nil
	}

	originalUpdateEvent := handlers.UpdateEvent
	defer func() {
		handlers.UpdateEvent = originalUpdateEvent
	}()
	handlers.UpdateEvent = func(event *models.Event) error {
		event.ID = 1
		return errors.New("simulate error update event handler")
	}

	req, _ := http.NewRequest(
		http.MethodPut,
		strings.Replace(UPDATE_EVENT_PATH, ":eventId", "1", 1),
		body,
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "sample-token-userid-1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Could not update event", response["message"])
	assert.NotEmpty(t, response["error"].(string))
}

func TestDeleteEvent_Success(t *testing.T) {
	router := setupRouter()

	// mock utils.VerifyToken
	originalVerifyToken := utils.VerifyToken
	defer func() {
		utils.VerifyToken = originalVerifyToken
	}()
	utils.VerifyToken = func(token string) (int64, error) {
		return 1, nil
	}

	// mock handlers.GetEventByID
	originalGetEventByID := handlers.GetEventByID
	defer func() {
		handlers.GetEventByID = originalGetEventByID
	}()
	handlers.GetEventByID = func(eventId int64) (*models.Event, error) {
		return &models.Event{
				ID:          1,
				Name:        "Go Workshop Jakarta",
				Description: "A beginner-friendly workshop covering Go fundamentals and best practices.",
				Location:    "Jakarta",
				DateTime:    time.Date(2025, 12, 16, 9, 0, 0, 0, time.FixedZone("WIB", 7*3600)),
				UserID:      1,
			},
			nil
	}

	// mock handlers.DeleteEvent
	originalDeleteEvent := handlers.DeleteEvent
	defer func() {
		handlers.DeleteEvent = originalDeleteEvent
	}()
	handlers.DeleteEvent = func(event *models.Event) error {
		return nil
	}

	// simulate hit API
	req, _ := http.NewRequest(
		http.MethodDelete,
		strings.Replace(DELETE_EVENT_PATH, ":eventId", "1", 1),
		http.NoBody,
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "sample-token-userid-1")

	// using w to simulate expected response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// expected response statusCode
	assert.Equal(t, http.StatusOK, w.Code)

	// convert JSON to []byte
	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// expected response data
	assert.Equal(t, "Event deleted successfully", response["message"])
}

func TestDeleteEvent_ErrorParseEventId(t *testing.T) {
	router := setupRouter()

	// mock utils.VerifyToken
	originalVerifyToken := utils.VerifyToken
	defer func() {
		utils.VerifyToken = originalVerifyToken
	}()
	utils.VerifyToken = func(token string) (int64, error) {
		return 1, nil
	}

	// simulate hit API
	req, _ := http.NewRequest(
		http.MethodDelete,
		strings.Replace(DELETE_EVENT_PATH, ":eventId", "invalidEventId", 1),
		http.NoBody,
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "sample-token-userid-1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Invalid event ID", response["message"])
	assert.NotEmpty(t, response["error"].(string))
}

func TestDeleteEvent_ErrorGetEventById(t *testing.T) {
	router := setupRouter()

	// mock utils.VerifyToken
	originalVerifyToken := utils.VerifyToken
	defer func() {
		utils.VerifyToken = originalVerifyToken
	}()
	utils.VerifyToken = func(token string) (int64, error) {
		return 1, nil
	}

	// mock handlers.GetEventByID
	originalGetEventByID := handlers.GetEventByID
	defer func() {
		handlers.GetEventByID = originalGetEventByID
	}()
	handlers.GetEventByID = func(eventId int64) (*models.Event, error) {
		return &models.Event{},
			errors.New("simulate error get event by id")
	}

	// simulate hit API
	req, _ := http.NewRequest(
		http.MethodDelete,
		strings.Replace(DELETE_EVENT_PATH, ":eventId", "1", 1),
		http.NoBody,
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "sample-token-userid-1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Event not found", response["message"])
	assert.NotEmpty(t, response["error"].(string))
}

func TestDeleteEvent_ErrorUnauthorize(t *testing.T) {
	router := setupRouter()

	// mock utils.VerifyToken
	originalVerifyToken := utils.VerifyToken
	defer func() {
		utils.VerifyToken = originalVerifyToken
	}()
	utils.VerifyToken = func(token string) (int64, error) {
		return 67, nil
	}

	// mock handlers.GetEventByID
	originalGetEventByID := handlers.GetEventByID
	defer func() {
		handlers.GetEventByID = originalGetEventByID
	}()
	handlers.GetEventByID = func(eventId int64) (*models.Event, error) {
		return &models.Event{
				ID:          1,
				Name:        "Go Workshop Jakarta",
				Description: "A beginner-friendly workshop covering Go fundamentals and best practices.",
				Location:    "Jakarta",
				DateTime:    time.Date(2025, 12, 16, 9, 0, 0, 0, time.FixedZone("WIB", 7*3600)),
				UserID:      1,
			},
			nil
	}

	// simulate hit API
	req, _ := http.NewRequest(
		http.MethodDelete,
		strings.Replace(DELETE_EVENT_PATH, ":eventId", "1", 1),
		http.NoBody,
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "sample-token-userid-67")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Not authorized to delete event", response["message"])
	assert.Equal(t, "not authorized to delete event", response["error"])
}

func TestDeleteEvent_ErrorDeleteEventHandler(t *testing.T) {
	router := setupRouter()

	// mock utils.VerifyToken
	originalVerifyToken := utils.VerifyToken
	defer func() {
		utils.VerifyToken = originalVerifyToken
	}()
	utils.VerifyToken = func(token string) (int64, error) {
		return 1, nil
	}

	// mock handlers.GetEventByID
	originalGetEventByID := handlers.GetEventByID
	defer func() {
		handlers.GetEventByID = originalGetEventByID
	}()
	handlers.GetEventByID = func(eventId int64) (*models.Event, error) {
		return &models.Event{
				ID:          1,
				Name:        "Go Workshop Jakarta",
				Description: "A beginner-friendly workshop covering Go fundamentals and best practices.",
				Location:    "Jakarta",
				DateTime:    time.Date(2025, 12, 16, 9, 0, 0, 0, time.FixedZone("WIB", 7*3600)),
				UserID:      1,
			},
			nil
	}

	// mock handlers.DeleteEvent
	originalDeleteEvent := handlers.DeleteEvent
	defer func() {
		handlers.DeleteEvent = originalDeleteEvent
	}()
	handlers.DeleteEvent = func(event *models.Event) error {
		return errors.New("simulate error delete event handler")
	}

	// simulate hit API
	req, _ := http.NewRequest(
		http.MethodDelete,
		strings.Replace(DELETE_EVENT_PATH, ":eventId", "1", 1),
		http.NoBody,
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "sample-token-userid-1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Could not delete event", response["message"])
	assert.NotEmpty(t, response["error"])
}
