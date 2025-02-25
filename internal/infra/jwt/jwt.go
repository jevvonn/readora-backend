package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jevvonn/readora-backend/config"
)

func CreateAuthToken(userId string, username string, role string) (string, error) {
	data := jwt.MapClaims{
		"sub":      userId,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
		"username": username,
		"role":     role,
	}

	config := config.Load()
	key := []byte(config.JWTSecret)

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	tokenString, err := t.SignedString(key)

	return tokenString, err
}

func ParseAuthToken(tokenString string) (jwt.MapClaims, error) {
	config := config.Load()
	key := []byte(config.JWTSecret)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return key, nil
	})

	if err != nil {
		return jwt.MapClaims{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	} else {
		return jwt.MapClaims{}, err
	}
}
