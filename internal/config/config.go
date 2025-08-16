package config

import (
	"fmt"
	"log"
	"os"
)

// Configuration constants
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

// Config structure
type Config struct {
	Port         string
	Host         string
	ProjectsFile string
	TasksFile    string
}

// Get default configuration
func GetDefaultConfig() *Config {
	return &Config{
		Port:         DefaultPort,
		Host:         DefaultHost,
		ProjectsFile: ProjectsFile,
		TasksFile:    TasksFile,
	}
}

// Load configuration from environment variables
func LoadFromEnv() *Config {
	config := GetDefaultConfig()

	// Check environment variables
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

// Print configuration
func PrintConfig(config *Config) {
	log.Printf("Configuration:")
	log.Printf("  Port: %s", config.Port)
	log.Printf("  Host: %s", config.Host)
	log.Printf("  Projects file: %s", config.ProjectsFile)
	log.Printf("  Tasks file: %s", config.TasksFile)
}

// Validate configuration
func ValidateConfig(config *Config) error {
	if config.Port == "" {
		return fmt.Errorf("Port cannot be empty")
	}
	if config.Host == "" {
		return fmt.Errorf("Host cannot be empty")
	}
	return nil
}
