package main

import (
	"log"

	"github.com/jevvonn/reodora-backend/internal/infra"
)

// First Time Init For Swagger Documentation
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description "Type 'Bearer TOKEN' to correctly set the API Key"

func main() {
	err := infra.Bootstrap()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
