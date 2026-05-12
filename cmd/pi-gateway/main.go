package main

import (
	"log"

	"github.com/zmoony/pi-gateway/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatalf("bootstrap application: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("run application: %v", err)
	}
}
