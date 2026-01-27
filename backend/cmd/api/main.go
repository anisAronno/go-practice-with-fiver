package main

import (
	"gofiver/internal/bootstrap"
	"log"
)

func main() {
	app, err := bootstrap.NewApplication()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	if err := app.RunSeeders(); err != nil {
		log.Printf("Seeder warning: %v", err)
	}

	if err := app.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
