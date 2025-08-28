package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"Projekt_go/internal/config"
	"Projekt_go/internal/models"
)

var (
	cfg *config.Config
)

// Initialize storage with configuration
func Initialize(config *config.Config) {
	cfg = config
}

// Load projects from file
func LoadProjects() error {
	file, err := os.Open(cfg.ProjectsFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Projects file doesn't exist, creating new one")
			return SaveProjects()
		}
		return err
	}
	defer file.Close()

	var list []models.Project
	err = json.NewDecoder(file).Decode(&list)
	if err != nil {
		return err
	}

	// Clear existing projects and reload
	projects := make(map[int]models.Project)
	maxID := 0
	for _, p := range list {
		projects[p.ID] = p
		if p.ID > maxID {
			maxID = p.ID
		}
	}

	// Update the global projects map
	models.ResetIDCounters()
	for _, p := range projects {
		models.SaveProject(p)
	}

	return nil
}

// Save projects to file
func SaveProjects() error {
	projects := models.GetAllProjects()

	data, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cfg.ProjectsFile, data, 0644)
}

// Load tasks from file
func LoadTasks() error {
	file, err := os.Open(cfg.TasksFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Tasks file doesn't exist, creating new one")
			return SaveTasks()
		}
		return err
	}
	defer file.Close()

	var list []models.Task
	err = json.NewDecoder(file).Decode(&list)
	if err != nil {
		return err
	}

	// Clear existing tasks and reload
	tasks := make(map[int]models.Task)
	maxID := 0
	for _, t := range list {
		tasks[t.ID] = t
		if t.ID > maxID {
			maxID = t.ID
		}
	}

	// Update the global tasks map
	models.ResetIDCounters()
	for _, t := range tasks {
		models.SaveTask(t)
	}

	return nil
}

// Save tasks to file
func SaveTasks() error {
	tasks := models.GetAllTasks()

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cfg.TasksFile, data, 0644)
}

// Export data to JSON
func ExportToJSON() ([]byte, error) {
	data := map[string]interface{}{
		"projects": models.GetAllProjects(),
		"tasks":    models.GetAllTasks(),
		"stats": map[string]interface{}{
			"total_projects": models.GetProjectsCount(),
			"total_tasks":    models.GetTasksCount(),
			"completed_tasks": func() int {
				tasks := models.GetAllTasks()
				count := 0
				for _, task := range tasks {
					if task.Done {
						count++
					}
				}
				return count
			}(),
		},
	}

	return json.MarshalIndent(data, "", "  ")
}

// Export data to CSV
func ExportToCSV() (string, error) {
	projects := models.GetAllProjects()
	tasks := models.GetAllTasks()

	csv := "ID,Name,Description,Status,CreatedAt,UpdatedAt\n"
	for _, p := range projects {
		csv += fmt.Sprintf("%d,%s,%s,%s,%s,%s\n",
			p.ID, p.Name, p.Description, p.Status,
			p.CreatedAt.Format("2006-01-02 15:04:05"),
			p.UpdatedAt.Format("2006-01-02 15:04:05"))
	}

	csv += "\nTaskID,ProjectID,Title,Description,Priority,Done,Deadline,CreatedAt,UpdatedAt\n"
	for _, t := range tasks {
		deadline := ""
		if t.Deadline != nil {
			deadline = t.Deadline.Format("2006-01-02 15:04:05")
		}
		csv += fmt.Sprintf("%d,%d,%s,%s,%d,%t,%s,%s,%s\n",
			t.ID, t.ProjectID, t.Title, t.Description, t.Priority, t.Done,
			deadline,
			t.CreatedAt.Format("2006-01-02 15:04:05"),
			t.UpdatedAt.Format("2006-01-02 15:04:05"))
	}

	return csv, nil
}
