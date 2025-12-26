package routes

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/event/handlers"
	"example.com/event/models"
	"github.com/stretchr/testify/assert"
)

var signUpPayload = map[string]string{
	"email":    "test@example.com",
	"password": "password123",
}

var loginPayload = map[string]string{
	"email":    "test@example.com",
	"password": "password123",
}

func TestSignUp_Success(t *testing.T) {
	router := setupRouter()

	body := toJSON(t, signUpPayload)

	// keep original handler to restore global state after test
	originalSaveUser := handlers.SaveUser

	// restore original handler after test to avoid side effects
	defer func() {
		handlers.SaveUser = originalSaveUser
	}()

	// mock handlers.SaveUser
	handlers.SaveUser = func(user *models.User) error {
		user.ID = 1 // simulate DB insert
		return nil
	}

	// simulate hit API
	req, _ := http.NewRequest(
		http.MethodPost,
		SIGNUP_PATH,
		body,
	)
	req.Header.Set("Content-Type", "application/json")

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
	assert.Equal(t, "User created successfully", response["message"])
	assert.Equal(t, "test@example.com", data["email"])
	assert.Equal(t, float64(1), data["id"]) // JSON numbers → float64
}

func TestSignUp_ErrorShouldBindJSON(t *testing.T) {
	router := setupRouter()

	// send invalid JSON to trigger ShouldBindJSON error
	req, _ := http.NewRequest(
		http.MethodPost,
		"/signup",
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

func TestSignUp_ErrorSaveUserHandler(t *testing.T) {
	router := setupRouter()

	body := toJSON(t, signUpPayload)

	originalSaveUser := handlers.SaveUser
	defer func() {
		handlers.SaveUser = originalSaveUser
	}()

	handlers.SaveUser = func(user *models.User) error {
		return errors.New("simulate error handlers.SaveUser")
	}

	req, _ := http.NewRequest(
		http.MethodPost,
		SIGNUP_PATH,
		body,
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Could not create user", response["message"])
	assert.NotEmpty(t, response["error"].(string))
}

func TestLogin_Success(t *testing.T) {
	router := setupRouter()

	body := toJSON(t, loginPayload)

	// keep original handler to restore global state after test
	originalValidateCredentials := handlers.ValidateCredentials

	// restore original handler after test to avoid side effects
	defer func() {
		handlers.ValidateCredentials = originalValidateCredentials
	}()

	// mock handlers.SaveUser
	handlers.ValidateCredentials = func(user *models.User) error {
		return nil
	}

	originalGenerateToken := handlers.GenerateToken
	defer func() {
		handlers.GenerateToken = originalGenerateToken
	}()
	handlers.GenerateToken = func(email string, userId int64) (string, error) {
		return "token", nil
	}

	// simulate hit API
	req, _ := http.NewRequest(
		http.MethodPost,
		LOGIN_PATH,
		body,
	)
	req.Header.Set("Content-Type", "application/json")

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
	assert.Equal(t, "Login successful", response["message"])
	assert.NotEmpty(t, response["token"].(string))
}

func TestLogin_ErrorShouldBindJSON(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest(
		http.MethodPost,
		LOGIN_PATH,
		bytes.NewBufferString("simulate invalid JSON"),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Could not parse request", response["message"])
	assert.NotEmpty(t, response["error"].(string))
}

func TestLogin_ErrorValidateCredentials(t *testing.T) {
	router := setupRouter()

	body := toJSON(t, loginPayload)

	originalValidateCredentials := handlers.ValidateCredentials

	defer func() {
		handlers.ValidateCredentials = originalValidateCredentials
	}()

	handlers.ValidateCredentials = func(user *models.User) error {
		return errors.New("simulate error handlers.ValidateCredentials")
	}

	req, _ := http.NewRequest(
		http.MethodPost,
		LOGIN_PATH,
		body,
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Could not authenticate user", response["message"])
	assert.NotEmpty(t, response["error"].(string))
}

func TestLogin_ErrorGenerateToken(t *testing.T) {
	router := setupRouter()

	body := toJSON(t, loginPayload)

	originalValidateCredentials := handlers.ValidateCredentials
	defer func() {
		handlers.ValidateCredentials = originalValidateCredentials
	}()
	handlers.ValidateCredentials = func(user *models.User) error {
		return nil
	}

	originalGenerateToken := handlers.GenerateToken
	defer func() {
		handlers.GenerateToken = originalGenerateToken
	}()
	handlers.GenerateToken = func(email string, userId int64) (string, error) {
		return "", errors.New("simulate error handlers.GenerateToken")
	}

	req, _ := http.NewRequest(
		http.MethodPost,
		LOGIN_PATH,
		body,
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Could not generate token", response["message"])
	assert.NotEmpty(t, response["error"].(string))
}
