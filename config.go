package main

import (
	"fmt"
	"log"
	"os"
)

// konfiguracja aplikacji

// stałe konfiguracyjne
const (
	DefaultPort    = "8080"
	DefaultHost    = "localhost"
	ProjectsFile   = "data_projects.json"
	TasksFile      = "data_tasks.json"
	MaxProjectName = 100
	MaxTaskTitle   = 200
	MaxDescription = 1000
	MaxPriority    = 5
	MinPriority    = 0
)

// struktura konfiguracji
type Config struct {
	Port         string
	Host         string
	ProjectsFile string
	TasksFile    string
}

// domyślna konfiguracja
func getDefaultConfig() *Config {
	return &Config{
		Port:         DefaultPort,
		Host:         DefaultHost,
		ProjectsFile: ProjectsFile,
		TasksFile:    TasksFile,
	}
}

// wczytuje konfigurację z zmiennych środowiskowych
func loadConfigFromEnv() *Config {
	config := getDefaultConfig()

	// sprawdzamy zmienne środowiskowe
	if port := os.Getenv("PORT"); port != "" {
		config.Port = port
	}

	if host := os.Getenv("HOST"); host != "" {
		config.Host = host
	}

	if projectsFile := os.Getenv("PROJECTS_FILE"); projectsFile != "" {
		config.ProjectsFile = projectsFile
	}

	if tasksFile := os.Getenv("TASKS_FILE"); tasksFile != "" {
		config.TasksFile = tasksFile
	}

	return config
}

// wyświetla konfigurację
func printConfig(config *Config) {
	log.Printf("Konfiguracja:")
	log.Printf("  Port: %s", config.Port)
	log.Printf("  Host: %s", config.Host)
	log.Printf("  Plik projektów: %s", config.ProjectsFile)
	log.Printf("  Plik zadań: %s", config.TasksFile)
}

// sprawdza czy konfiguracja jest poprawna
func validateConfig(config *Config) error {
	if config.Port == "" {
		return fmt.Errorf("Port cannot be empty")
	}
	if config.Host == "" {
		return fmt.Errorf("Host cannot be empty")
	}
	return nil
}
