package main

import "fmt"

func validateProject(p *Project) error {
	if p.Name == "" {
		return fmt.Errorf("Project name is required")
	}
	if len(p.Name) > 100 {
		return fmt.Errorf("Project name too long (max 100 characters)")
	}
	return nil
}

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
	// Deadline jest opcjonalny, więc nie walidujemy go
	return nil
}
