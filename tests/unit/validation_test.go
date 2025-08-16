package unit

import (
	"testing"
	"time"

	"Projekt_go/internal/config"
	"Projekt_go/internal/models"
	"Projekt_go/internal/validation"
)

func TestValidateProject(t *testing.T) {
	// Test valid project
	validProject := &models.Project{
		Name:        "Test Project",
		Description: "Test Description",
		Status:      "active",
	}

	err := validation.ValidateProject(validProject)
	if err != nil {
		t.Errorf("Expected valid project to pass validation, got error: %v", err)
	}

	// Test empty name
	invalidProject := &models.Project{
		Name:        "",
		Description: "Test Description",
	}

	err = validation.ValidateProject(invalidProject)
	if err == nil {
		t.Error("Expected project with empty name to fail validation")
	}

	// Test name too long
	longName := ""
	for i := 0; i < config.MaxProjectName+1; i++ {
		longName += "a"
	}
	invalidProject.Name = longName

	err = validation.ValidateProject(invalidProject)
	if err == nil {
		t.Error("Expected project with long name to fail validation")
	}

	// Test invalid status
	invalidProject.Name = "Test Project"
	invalidProject.Status = "invalid_status"

	err = validation.ValidateProject(invalidProject)
	if err == nil {
		t.Error("Expected project with invalid status to fail validation")
	}
}

func TestValidateTask(t *testing.T) {
	// Test valid task
	validTask := &models.Task{
		Title:          "Test Task",
		Description:    "Test Description",
		Priority:       3,
		EstimatedHours: 5.0,
		ActualHours:    4.0,
	}

	err := validation.ValidateTask(validTask)
	if err != nil {
		t.Errorf("Expected valid task to pass validation, got error: %v", err)
	}

	// Test empty title
	invalidTask := &models.Task{
		Title:       "",
		Description: "Test Description",
		Priority:    1,
	}

	err = validation.ValidateTask(invalidTask)
	if err == nil {
		t.Error("Expected task with empty title to fail validation")
	}

	// Test title too long
	longTitle := ""
	for i := 0; i < config.MaxTaskTitle+1; i++ {
		longTitle += "a"
	}
	invalidTask.Title = longTitle

	err = validation.ValidateTask(invalidTask)
	if err == nil {
		t.Error("Expected task with long title to fail validation")
	}

	// Test invalid priority
	invalidTask.Title = "Test Task"
	invalidTask.Priority = 10

	err = validation.ValidateTask(invalidTask)
	if err == nil {
		t.Error("Expected task with invalid priority to fail validation")
	}

	// Test negative estimated hours
	invalidTask.Priority = 1
	invalidTask.EstimatedHours = -1.0

	err = validation.ValidateTask(invalidTask)
	if err == nil {
		t.Error("Expected task with negative estimated hours to fail validation")
	}

	// Test negative actual hours
	invalidTask.EstimatedHours = 1.0
	invalidTask.ActualHours = -1.0

	err = validation.ValidateTask(invalidTask)
	if err == nil {
		t.Error("Expected task with negative actual hours to fail validation")
	}
}

func TestValidateBadge(t *testing.T) {
	// Test valid badge
	validBadge := &models.Badge{
		Text:       "Test Badge",
		Color:      "#ffffff",
		Background: "#000000",
		Type:       "status",
	}

	err := validation.ValidateBadge(validBadge)
	if err != nil {
		t.Errorf("Expected valid badge to pass validation, got error: %v", err)
	}

	// Test empty text
	invalidBadge := &models.Badge{
		Text:       "",
		Color:      "#ffffff",
		Background: "#000000",
	}

	err = validation.ValidateBadge(invalidBadge)
	if err == nil {
		t.Error("Expected badge with empty text to fail validation")
	}

	// Test text too long
	longText := ""
	for i := 0; i < 51; i++ {
		longText += "a"
	}
	invalidBadge.Text = longText

	err = validation.ValidateBadge(invalidBadge)
	if err == nil {
		t.Error("Expected badge with long text to fail validation")
	}

	// Test empty color
	invalidBadge.Text = "Test Badge"
	invalidBadge.Color = ""

	err = validation.ValidateBadge(invalidBadge)
	if err == nil {
		t.Error("Expected badge with empty color to fail validation")
	}

	// Test empty background
	invalidBadge.Color = "#ffffff"
	invalidBadge.Background = ""

	err = validation.ValidateBadge(invalidBadge)
	if err == nil {
		t.Error("Expected badge with empty background to fail validation")
	}

	// Test invalid type
	invalidBadge.Background = "#000000"
	invalidBadge.Type = "invalid_type"

	err = validation.ValidateBadge(invalidBadge)
	if err == nil {
		t.Error("Expected badge with invalid type to fail validation")
	}
}

func TestValidateID(t *testing.T) {
	// Test valid ID
	err := validation.ValidateID(1)
	if err != nil {
		t.Errorf("Expected valid ID to pass validation, got error: %v", err)
	}

	// Test zero ID
	err = validation.ValidateID(0)
	if err == nil {
		t.Error("Expected zero ID to fail validation")
	}

	// Test negative ID
	err = validation.ValidateID(-1)
	if err == nil {
		t.Error("Expected negative ID to fail validation")
	}
}

func TestValidateStringLength(t *testing.T) {
	// Test valid string
	err := validation.ValidateStringLength("test", 10)
	if err != nil {
		t.Errorf("Expected valid string to pass validation, got error: %v", err)
	}

	// Test string too long
	err = validation.ValidateStringLength("test", 3)
	if err == nil {
		t.Error("Expected long string to fail validation")
	}
}

func TestValidateStringNotEmpty(t *testing.T) {
	// Test non-empty string
	err := validation.ValidateStringNotEmpty("test", "field")
	if err != nil {
		t.Errorf("Expected non-empty string to pass validation, got error: %v", err)
	}

	// Test empty string
	err = validation.ValidateStringNotEmpty("", "field")
	if err == nil {
		t.Error("Expected empty string to fail validation")
	}
}

func TestValidateNumberRange(t *testing.T) {
	// Test valid number
	err := validation.ValidateNumberRange(5, 1, 10, "field")
	if err != nil {
		t.Errorf("Expected valid number to pass validation, got error: %v", err)
	}

	// Test number too low
	err = validation.ValidateNumberRange(0, 1, 10, "field")
	if err == nil {
		t.Error("Expected low number to fail validation")
	}

	// Test number too high
	err = validation.ValidateNumberRange(11, 1, 10, "field")
	if err == nil {
		t.Error("Expected high number to fail validation")
	}
}

func TestValidateArrayNotEmpty(t *testing.T) {
	// Test non-empty array
	err := validation.ValidateArrayNotEmpty([]string{"test"}, "field")
	if err != nil {
		t.Errorf("Expected non-empty array to pass validation, got error: %v", err)
	}

	// Test empty array
	err = validation.ValidateArrayNotEmpty([]string{}, "field")
	if err == nil {
		t.Error("Expected empty array to fail validation")
	}
}

func TestValidateFutureDate(t *testing.T) {
	// Test future date
	futureDate := time.Now().Add(24 * time.Hour)
	err := validation.ValidateFutureDate(futureDate, "field")
	if err != nil {
		t.Errorf("Expected future date to pass validation, got error: %v", err)
	}

	// Test past date
	pastDate := time.Now().Add(-24 * time.Hour)
	err = validation.ValidateFutureDate(pastDate, "field")
	if err == nil {
		t.Error("Expected past date to fail validation")
	}
}

func TestValidatePastDate(t *testing.T) {
	// Test past date
	pastDate := time.Now().Add(-24 * time.Hour)
	err := validation.ValidatePastDate(pastDate, "field")
	if err != nil {
		t.Errorf("Expected past date to pass validation, got error: %v", err)
	}

	// Test future date
	futureDate := time.Now().Add(24 * time.Hour)
	err = validation.ValidatePastDate(futureDate, "field")
	if err == nil {
		t.Error("Expected future date to fail validation")
	}
}

func TestValidateConfig(t *testing.T) {
	// Test valid config
	validConfig := &config.Config{
		Port: "8080",
		Host: "localhost",
	}

	err := validation.ValidateConfig(validConfig)
	if err != nil {
		t.Errorf("Expected valid config to pass validation, got error: %v", err)
	}

	// Test empty port
	invalidConfig := &config.Config{
		Port: "",
		Host: "localhost",
	}

	err = validation.ValidateConfig(invalidConfig)
	if err == nil {
		t.Error("Expected config with empty port to fail validation")
	}

	// Test empty host
	invalidConfig.Port = "8080"
	invalidConfig.Host = ""

	err = validation.ValidateConfig(invalidConfig)
	if err == nil {
		t.Error("Expected config with empty host to fail validation")
	}
}
