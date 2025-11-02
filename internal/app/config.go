package app

import (
	"fmt"
	"log"
	"os"
)

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

// Config przechowuje ustawienia aplikacji.
type Config struct {
	Port         string
	Host         string
	ProjectsFile string
	TasksFile    string
}

var currentConfig = defaultConfig()

func defaultConfig() *Config {
	return &Config{
		Port:         DefaultPort,
		Host:         DefaultHost,
		ProjectsFile: ProjectsFile,
		TasksFile:    TasksFile,
	}
}

// SetCurrentConfig ustawia aktualną konfigurację aplikacji.
func SetCurrentConfig(cfg *Config) {
	if cfg == nil {
		currentConfig = defaultConfig()
		return
	}
	currentConfig = cfg
}

// GetCurrentConfig zwraca aktualnie obowiązującą konfigurację.
func GetCurrentConfig() *Config {
	if currentConfig == nil {
		return defaultConfig()
	}
	return currentConfig
}

// LoadConfigFromEnv pobiera konfigurację ze zmiennych środowiskowych.
func LoadConfigFromEnv() *Config {
	cfg := defaultConfig()

	if port := os.Getenv("PORT"); port != "" {
		cfg.Port = port
	}
	if host := os.Getenv("HOST"); host != "" {
		cfg.Host = host
	}

	if projectsFile := os.Getenv("PROJECTS_FILE"); projectsFile != "" {
		cfg.ProjectsFile = projectsFile
	}

	if tasksFile := os.Getenv("TASKS_FILE"); tasksFile != "" {
		cfg.TasksFile = tasksFile
	}

	return cfg
}

// PrintConfig wypisuje konfigurację do logów.
func PrintConfig(cfg *Config) {
	log.Printf("Konfiguracja:")
	log.Printf("  Port: %s", cfg.Port)
	log.Printf("  Host: %s", cfg.Host)
	log.Printf("  Plik projektów: %s", cfg.ProjectsFile)
	log.Printf("  Plik zadań: %s", cfg.TasksFile)
}

// ValidateConfig sprawdza poprawność konfiguracji.
func ValidateConfig(cfg *Config) error {
	if cfg.Port == "" {
		return fmt.Errorf("Port cannot be empty")
	}
	if cfg.Host == "" {
		return fmt.Errorf("Host cannot be empty")
	}
	return nil
}
