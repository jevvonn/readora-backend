package main

import (
	"log"

	"github.com/jevvonn/reodora-backend/internal/infra"
)

func main() {
	err := infra.Bootstrap()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
