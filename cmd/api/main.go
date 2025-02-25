package main

import (
	"log"

	"github.com/jevvonn/readora-backend/internal/bootstrap"
)

// First Time Init For Swagger Documentation
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description "Type 'Bearer TOKEN' to correctly set the API Key"

func main() {
	err := bootstrap.Start()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
