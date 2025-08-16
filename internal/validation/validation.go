package validation

import (
	"fmt"
	"time"

	"Projekt_go/internal/config"
	"Projekt_go/internal/models"
)

// Validate project
func ValidateProject(p *models.Project) error {
	if p.Name == "" {
		return fmt.Errorf("Project name is required")
	}
	if len(p.Name) > config.MaxProjectName {
		return fmt.Errorf("Project name too long (max %d characters)", config.MaxProjectName)
	}
	if len(p.Description) > config.MaxDescription {
		return fmt.Errorf("Project description too long (max %d characters)", config.MaxDescription)
	}
	if p.Status != "" && p.Status != "active" && p.Status != "completed" && p.Status != "archived" {
		return fmt.Errorf("Invalid project status")
	}
	if p.Badge != nil {
		if err := ValidateBadge(p.Badge); err != nil {
			return fmt.Errorf("Invalid project badge: %v", err)
		}
	}
	return nil
}

// Validate task
func ValidateTask(t *models.Task) error {
	if t.Title == "" {
		return fmt.Errorf("Task title is required")
	}
	if len(t.Title) > config.MaxTaskTitle {
		return fmt.Errorf("Task title too long (max %d characters)", config.MaxTaskTitle)
	}
	if len(t.Description) > config.MaxDescription {
		return fmt.Errorf("Task description too long (max %d characters)", config.MaxDescription)
	}
	if t.Priority < config.MinPriority || t.Priority > config.MaxPriority {
		return fmt.Errorf("Priority must be between %d and %d", config.MinPriority, config.MaxPriority)
	}
	if t.EstimatedHours < 0 {
		return fmt.Errorf("Estimated hours cannot be negative")
	}
	if t.ActualHours < 0 {
		return fmt.Errorf("Actual hours cannot be negative")
	}
	if t.Badge != nil {
		if err := ValidateBadge(t.Badge); err != nil {
			return fmt.Errorf("Invalid task badge: %v", err)
		}
	}
	// deadline is optional
	return nil
}

// Validate badge
func ValidateBadge(b *models.Badge) error {
	if b.Text == "" {
		return fmt.Errorf("Badge text is required")
	}
	if len(b.Text) > 50 {
		return fmt.Errorf("Badge text too long (max 50 characters)")
	}
	if b.Color == "" {
		return fmt.Errorf("Badge color is required")
	}
	if b.Background == "" {
		return fmt.Errorf("Badge background is required")
	}
	if b.Type != "" && b.Type != "status" && b.Type != "priority" && b.Type != "category" && b.Type != "custom" {
		return fmt.Errorf("Invalid badge type")
	}
	return nil
}

// Validate ID
func ValidateID(id int) error {
	if id <= 0 {
		return fmt.Errorf("ID must be positive")
	}
	return nil
}

// Validate string length
func ValidateStringLength(s string, maxLength int) error {
	if len(s) > maxLength {
		return fmt.Errorf("String too long (max %d characters)", maxLength)
	}
	return nil
}

// Validate string is not empty
func ValidateStringNotEmpty(s string, fieldName string) error {
	if s == "" {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	return nil
}

// Validate number is in range
func ValidateNumberRange(num, min, max int, fieldName string) error {
	if num < min || num > max {
		return fmt.Errorf("%s must be between %d and %d", fieldName, min, max)
	}
	return nil
}

// Validate array is not empty
func ValidateArrayNotEmpty(arr []string, fieldName string) error {
	if len(arr) == 0 {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	return nil
}

// Validate date is in the future
func ValidateFutureDate(date time.Time, fieldName string) error {
	if date.Before(time.Now()) {
		return fmt.Errorf("%s must be in the future", fieldName)
	}
	return nil
}

// Validate date is in the past
func ValidatePastDate(date time.Time, fieldName string) error {
	if date.After(time.Now()) {
		return fmt.Errorf("%s must be in the past", fieldName)
	}
	return nil
}

// Validate configuration
func ValidateConfig(cfg *config.Config) error {
	if cfg.Port == "" {
		return fmt.Errorf("Port cannot be empty")
	}
	if cfg.Host == "" {
		return fmt.Errorf("Host cannot be empty")
	}
	return nil
}
