package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jevvonn/reodora-backend/config"
)

func CreateAuthToken(userId string, username string) (string, error) {
	data := jwt.MapClaims{
		"sub":      userId,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	}

	config := config.Load()
	key := []byte(config.JWTSecret)

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	tokenString, err := t.SignedString(key)

	return tokenString, err
}
