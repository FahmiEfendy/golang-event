package handlers

import (
	"example.com/event/models"
	"example.com/event/utils"
)

var SaveUser = func(user *models.User) error {
	return user.Save()
}

var ValidateCredentials = func(user *models.User) error {
	return user.ValidateCredentials()
}

var GenerateToken = func(email string, userId int64) (string, error) {
	return utils.GenerateToken(email, userId)
}
