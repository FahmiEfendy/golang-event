package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "EVENT-SECRET-KEY"

func GenerateToken(email string, userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":     email,
		"userId":    userId,
		"expiredAt": time.Now().Add(time.Hour * 2).Unix(), // expired at 2hrs
	})

	return token.SignedString([]byte(secretKey))
}
