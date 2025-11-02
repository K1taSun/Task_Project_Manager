package app

import "fmt"

func validateProject(p *Project) error {
	if p.Name == "" {
		return fmt.Errorf("project name is required")
	}
	if len(p.Name) > MaxProjectName {
		return fmt.Errorf("project name too long (max %d characters)", MaxProjectName)
	}
	if len(p.Description) > 500 {
		return fmt.Errorf("project description too long (max 500 characters)")
	}
	if p.Status != "" && p.Status != "active" && p.Status != "completed" && p.Status != "archived" {
		return fmt.Errorf("invalid project status")
	}
	if p.Badge != nil {
		if err := validateBadge(p.Badge); err != nil {
			return fmt.Errorf("invalid project badge: %v", err)
		}
	}
	return nil
}

func validateTask(t *Task) error {
	if t.Title == "" {
		return fmt.Errorf("task title is required")
	}
	if len(t.Title) > MaxTaskTitle {
		return fmt.Errorf("task title too long (max %d characters)", MaxTaskTitle)
	}
	if len(t.Description) > MaxDescription {
		return fmt.Errorf("task description too long (max %d characters)", MaxDescription)
	}
	if t.Priority < MinPriority || t.Priority > MaxPriority {
		return fmt.Errorf("priority must be between %d and %d", MinPriority, MaxPriority)
	}
	if t.EstimatedHours < 0 {
		return fmt.Errorf("estimated hours cannot be negative")
	}
	if t.ActualHours < 0 {
		return fmt.Errorf("actual hours cannot be negative")
	}
	if t.Badge != nil {
		if err := validateBadge(t.Badge); err != nil {
			return fmt.Errorf("invalid task badge: %v", err)
		}
	}
	return nil
}

func validateBadge(b *Badge) error {
	if b.Text == "" {
		return fmt.Errorf("badge text is required")
	}
	if len(b.Text) > 50 {
		return fmt.Errorf("badge text too long (max 50 characters)")
	}
	if b.Color == "" {
		return fmt.Errorf("badge color is required")
	}
	if b.Background == "" {
		return fmt.Errorf("badge background is required")
	}
	if b.Type != "" && b.Type != "status" && b.Type != "priority" && b.Type != "category" && b.Type != "custom" {
		return fmt.Errorf("invalid badge type")
	}
	return nil
}

func validateID(id int) error {
	if id <= 0 {
		return fmt.Errorf("id must be positive")
	}
	return nil
}
