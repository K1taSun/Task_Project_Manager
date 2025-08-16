package main

import (
	"log"
	"net/http"
	"os"

	"Projekt_go/internal/api"
	"Projekt_go/internal/config"
	"Projekt_go/internal/storage"
	"Projekt_go/internal/validation"
)

func main() {
	// Initialize logging
	log.Println("Starting Task Project Manager...")

	// Load configuration
	cfg := config.LoadFromEnv()
	config.PrintConfig(cfg)

	// Validate configuration
	err := validation.ValidateConfig(cfg)
	if err != nil {
		log.Fatal("Configuration error:", err)
	}

	// Initialize storage
	storage.Initialize(cfg)

	// Check if data files exist, create if not
	if _, err := os.Stat(cfg.ProjectsFile); os.IsNotExist(err) {
		log.Println("Creating file", cfg.ProjectsFile)
		err := storage.SaveProjects()
		if err != nil {
			log.Fatal("Error creating projects file:", err)
		}
	}

	if _, err := os.Stat(cfg.TasksFile); os.IsNotExist(err) {
		log.Println("Creating file", cfg.TasksFile)
		err := storage.SaveTasks()
		if err != nil {
			log.Fatal("Error creating tasks file:", err)
		}
	}

	// Load data
	err = storage.LoadProjects()
	if err != nil {
		log.Fatal("Error loading projects:", err)
	}
	err = storage.LoadTasks()
	if err != nil {
		log.Fatal("Error loading tasks:", err)
	}

	// Setup routes
	router := api.SetupRoutes()

	// Start server
	serverAddr := cfg.Host + ":" + cfg.Port
	log.Println("Server running on", serverAddr)
	log.Println("Open http://" + serverAddr + " in your browser")

	// Start server
	err = http.ListenAndServe(serverAddr, router)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
