package app

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// stałe konfiguracyjne
const (
	DefaultPort    = "8080"
	DefaultHost    = "localhost"
	ProjectsFile   = "data_projects.json"
	TasksFile      = "data_tasks.json"
	UsersFile      = "data_users.json"
	StaticDir      = "."
	MaxProjectName = 100
	MaxTaskTitle   = 200
	MaxDescription = 1000
	MaxPriority    = 5
	MinPriority    = 0
)

// Config przechowuje ustawienia aplikacji.
type Config struct {
	Port                  string
	Host                  string
	ProjectsFile          string
	TasksFile             string
	UsersFile             string
	StaticDir             string
	SessionSecret         string
	AdminRegistrationToken string
}

var (
	projectRootPath = findProjectRoot()
	currentConfig   = defaultConfig()
)

func defaultConfig() *Config {
	return &Config{
		Port:         DefaultPort,
		Host:         DefaultHost,
		ProjectsFile: filepath.Join(projectRootPath, ProjectsFile),
		TasksFile:    filepath.Join(projectRootPath, TasksFile),
		UsersFile:    filepath.Join(projectRootPath, UsersFile),
		StaticDir:    filepath.Join(projectRootPath, StaticDir),
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
		cfg.ProjectsFile = resolvePath(projectsFile)
	}

	if tasksFile := os.Getenv("TASKS_FILE"); tasksFile != "" {
		cfg.TasksFile = resolvePath(tasksFile)
	}

	if usersFile := os.Getenv("USERS_FILE"); usersFile != "" {
		cfg.UsersFile = resolvePath(usersFile)
	}

	if staticDir := os.Getenv("STATIC_DIR"); staticDir != "" {
		cfg.StaticDir = resolvePath(staticDir)
	}

	if secret := os.Getenv("SESSION_SECRET"); secret != "" {
		cfg.SessionSecret = secret
	}

	if token := os.Getenv("ADMIN_REGISTRATION_TOKEN"); token != "" {
		cfg.AdminRegistrationToken = token
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
	log.Printf("  Plik użytkowników: %s", cfg.UsersFile)
	log.Printf("  Katalog statyczny: %s", cfg.StaticDir)
	if cfg.AdminRegistrationToken != "" {
		log.Printf("  Token rejestracji admina: ustawiony")
	} else {
		log.Printf("  Token rejestracji admina: nie ustawiony (rejestracja admina wyłączona)")
	}
}

// ValidateConfig sprawdza poprawność konfiguracji.
func ValidateConfig(cfg *Config) error {
	if cfg.Port == "" {
		return fmt.Errorf("port cannot be empty")
	}
	if cfg.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}
	if cfg.UsersFile == "" {
		return fmt.Errorf("users file cannot be empty")
	}
	if cfg.StaticDir == "" {
		return fmt.Errorf("static dir cannot be empty")
	}
	if cfg.SessionSecret == "" {
		secret, err := generateRandomSecret()
		if err != nil {
			return fmt.Errorf("session secret cannot be empty and random generation failed: %v", err)
		}
		cfg.SessionSecret = secret
		log.Println("SESSION_SECRET not provided. Generated ephemeral secret for this run.")
	}
	return nil
}

func generateRandomSecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func resolvePath(p string) string {
	if p == "" {
		return ""
	}
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Join(projectRootPath, p)
}

func findProjectRoot() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	for {
		if _, err := os.Stat(filepath.Join(cwd, "go.mod")); err == nil {
			return cwd
		}
		parent := filepath.Dir(cwd)
		if parent == cwd {
			break
		}
		cwd = parent
	}
	return cwd
}
