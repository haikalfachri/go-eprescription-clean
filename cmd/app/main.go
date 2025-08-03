package main

import (
	"log"

	"go-eprescription-clean/config"
	"go-eprescription-clean/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
