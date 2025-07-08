package main

import "fmt"

func validateProject(p *Project) error {
	if p.Name == "" {
		return fmt.Errorf("Project name is required")
	}
	return nil
}

func validateTask(t *Task) error {
	if t.Title == "" {
		return fmt.Errorf("Task title is required")
	}
	if t.Priority < 0 || t.Priority > 5 {
		return fmt.Errorf("Priority must be between 0 and 5")
	}
	if t.Deadline.IsZero() {
		return fmt.Errorf("Deadline is required")
	}
	return nil
} 