package main

import (
	"fmt"
	"time"
)

// sprawdza czy projekt jest poprawny
func validateProject(p *Project) error {
	if p.Name == "" {
		return fmt.Errorf("Project name is required")
	}
	if len(p.Name) > 100 {
		return fmt.Errorf("Project name too long (max 100 characters)")
	}
	if len(p.Description) > 500 {
		return fmt.Errorf("Project description too long (max 500 characters)")
	}
	if p.Status != "" && p.Status != "active" && p.Status != "completed" && p.Status != "archived" {
		return fmt.Errorf("Invalid project status")
	}
	if p.Badge != nil {
		if err := validateBadge(p.Badge); err != nil {
			return fmt.Errorf("Invalid project badge: %v", err)
		}
	}
	return nil
}

// sprawdza czy zadanie jest poprawne
func validateTask(t *Task) error {
	if t.Title == "" {
		return fmt.Errorf("Task title is required")
	}
	if len(t.Title) > 200 {
		return fmt.Errorf("Task title too long (max 200 characters)")
	}
	if len(t.Description) > 1000 {
		return fmt.Errorf("Task description too long (max 1000 characters)")
	}
	if t.Priority < 0 || t.Priority > 5 {
		return fmt.Errorf("Priority must be between 0 and 5")
	}
	if t.EstimatedHours < 0 {
		return fmt.Errorf("Estimated hours cannot be negative")
	}
	if t.ActualHours < 0 {
		return fmt.Errorf("Actual hours cannot be negative")
	}
	if t.Badge != nil {
		if err := validateBadge(t.Badge); err != nil {
			return fmt.Errorf("Invalid task badge: %v", err)
		}
	}
	// deadline jest opcjonalny
	return nil
}

// sprawdza czy odznaka jest poprawna
func validateBadge(b *Badge) error {
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

// sprawdza czy ID jest poprawne
func validateID(id int) error {
	if id <= 0 {
		return fmt.Errorf("ID must be positive")
	}
	return nil
}

// sprawdza czy string ma odpowiednią długość
func validateStringLength(s string, maxLength int) error {
	if len(s) > maxLength {
		return fmt.Errorf("String too long (max %d characters)", maxLength)
	}
	return nil
}

// prosta funkcja do sprawdzania czy string nie jest pusty
func validateStringNotEmpty(s string, fieldName string) error {
	if s == "" {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	return nil
}

// prosta funkcja do sprawdzania czy liczba jest w zakresie
func validateNumberRange(num, min, max int, fieldName string) error {
	if num < min || num > max {
		return fmt.Errorf("%s must be between %d and %d", fieldName, min, max)
	}
	return nil
}

// prosta funkcja do sprawdzania czy tablica nie jest pusta
func validateArrayNotEmpty(arr []string, fieldName string) error {
	if len(arr) == 0 {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	return nil
}

// sprawdza czy data jest w przyszłości
func validateFutureDate(date time.Time, fieldName string) error {
	if date.Before(time.Now()) {
		return fmt.Errorf("%s must be in the future", fieldName)
	}
	return nil
}

// sprawdza czy data jest w przeszłości
func validatePastDate(date time.Time, fieldName string) error {
	if date.After(time.Now()) {
		return fmt.Errorf("%s must be in the past", fieldName)
	}
	return nil
}
