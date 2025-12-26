package routes

import (
	"net/http"

	"example.com/event/handlers"
	"example.com/event/models"
	"github.com/gin-gonic/gin"
)

type UserResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

func signUp(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse request",
			"error":   err.Error()})
		return
	}

	err = handlers.SaveUser(&user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not create user",
			"error":   err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"data": UserResponse{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}

func login(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse request",
			"error":   err.Error()})
		return
	}

	err = handlers.ValidateCredentials(&user)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Could not authenticate user",
			"error":   err.Error(),
		})
		return
	}

	token, err := handlers.GenerateToken(user.Email, user.ID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not generate token",
			"error":   err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}
