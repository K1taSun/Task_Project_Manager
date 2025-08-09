package main

import (
	"fmt"
	"testing"
)

// proste testy dla funkcji walidacji
// TODO: dodać więcej testów

func TestValidateProject(t *testing.T) {
	// test poprawnego projektu
	project := &Project{Name: "Test Project"}
	err := validateProject(project)
	if err != nil {
		t.Errorf("Project should be valid: %v", err)
	}

	// test pustej nazwy
	project.Name = ""
	err = validateProject(project)
	if err == nil {
		t.Error("Project with empty name should be invalid")
	}
}

func TestValidateTask(t *testing.T) {
	// test poprawnego zadania
	task := &Task{
		Title:       "Test Task",
		Description: "Test Description",
		Priority:    3,
	}
	err := validateTask(task)
	if err != nil {
		t.Errorf("Task should be valid: %v", err)
	}

	// test pustego tytułu
	task.Title = ""
	err = validateTask(task)
	if err == nil {
		t.Error("Task with empty title should be invalid")
	}
}

func TestValidateID(t *testing.T) {
	// test poprawnego ID
	err := validateID(1)
	if err != nil {
		t.Errorf("ID 1 should be valid: %v", err)
	}

	// test niepoprawnego ID
	err = validateID(0)
	if err == nil {
		t.Error("ID 0 should be invalid")
	}
}

// prosta funkcja do testowania
func TestMain(m *testing.M) {
	fmt.Println("Starting tests...")
	m.Run()
	fmt.Println("Tests finished.")
}
