package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"Projekt_go/internal/api"
	"Projekt_go/internal/config"
	"Projekt_go/internal/models"
	"Projekt_go/internal/storage"
)

func setupTestServer() *httptest.Server {
	// Initialize with test config
	cfg := &config.Config{
		Port:         "8080",
		Host:         "localhost",
		ProjectsFile: "test_projects.json",
		TasksFile:    "test_tasks.json",
	}

	storage.Initialize(cfg)
	models.ResetIDCounters()

	// Setup routes
	router := api.SetupRoutes()
	return httptest.NewServer(router)
}

func TestProjectsEndpoint(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Test GET /projects (empty)
	resp, err := http.Get(server.URL + "/projects")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var projects []models.Project
	err = json.NewDecoder(resp.Body).Decode(&projects)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(projects) != 0 {
		t.Errorf("Expected empty projects list, got %d projects", len(projects))
	}

	// Test POST /projects
	project := models.Project{
		Name:        "Test Project",
		Description: "Test Description",
		Status:      "active",
	}

	projectData, _ := json.Marshal(project)
	resp, err = http.Post(server.URL+"/projects", "application/json", bytes.NewBuffer(projectData))
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}

	var createdProject models.Project
	err = json.NewDecoder(resp.Body).Decode(&createdProject)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if createdProject.ID == 0 {
		t.Error("Expected project to have an ID")
	}
	if createdProject.Name != "Test Project" {
		t.Errorf("Expected project name 'Test Project', got '%s'", createdProject.Name)
	}

	// Test GET /projects (with data)
	resp, err = http.Get(server.URL + "/projects")
	if err != nil {
		t.Fatalf("Failed to get projects: %v", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&projects)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(projects) != 1 {
		t.Errorf("Expected 1 project, got %d", len(projects))
	}
}

func TestProjectEndpoint(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a project first
	project := models.Project{
		Name:        "Test Project",
		Description: "Test Description",
		Status:      "active",
	}

	projectData, _ := json.Marshal(project)
	resp, err := http.Post(server.URL+"/projects", "application/json", bytes.NewBuffer(projectData))
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}
	defer resp.Body.Close()

	var createdProject models.Project
	json.NewDecoder(resp.Body).Decode(&createdProject)

	// Test GET /projects/{id}
	resp, err = http.Get(server.URL + "/projects/" + string(rune(createdProject.ID)))
	if err != nil {
		t.Fatalf("Failed to get project: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var retrievedProject models.Project
	err = json.NewDecoder(resp.Body).Decode(&retrievedProject)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if retrievedProject.ID != createdProject.ID {
		t.Errorf("Expected project ID %d, got %d", createdProject.ID, retrievedProject.ID)
	}

	// Test PUT /projects/{id}
	updatedProject := models.Project{
		Name:        "Updated Project",
		Description: "Updated Description",
		Status:      "completed",
	}

	updatedData, _ := json.Marshal(updatedProject)
	req, _ := http.NewRequest("PUT", server.URL+"/projects/"+string(rune(createdProject.ID)), bytes.NewBuffer(updatedData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to update project: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Test DELETE /projects/{id}
	req, _ = http.NewRequest("DELETE", server.URL+"/projects/"+string(rune(createdProject.ID)), nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to delete project: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestTasksEndpoint(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a project first
	project := models.Project{
		Name:        "Test Project",
		Description: "Test Description",
		Status:      "active",
	}

	projectData, _ := json.Marshal(project)
	resp, err := http.Post(server.URL+"/projects", "application/json", bytes.NewBuffer(projectData))
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}
	defer resp.Body.Close()

	var createdProject models.Project
	json.NewDecoder(resp.Body).Decode(&createdProject)

	// Test POST /tasks
	task := models.Task{
		ProjectID:      createdProject.ID,
		Title:          "Test Task",
		Description:    "Test Description",
		Priority:       3,
		EstimatedHours: 5.0,
		Tags:           []string{"test", "integration"},
	}

	taskData, _ := json.Marshal(task)
	resp, err = http.Post(server.URL+"/tasks", "application/json", bytes.NewBuffer(taskData))
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}

	var createdTask models.Task
	err = json.NewDecoder(resp.Body).Decode(&createdTask)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if createdTask.ID == 0 {
		t.Error("Expected task to have an ID")
	}
	if createdTask.Title != "Test Task" {
		t.Errorf("Expected task title 'Test Task', got '%s'", createdTask.Title)
	}

	// Test GET /tasks with filtering
	resp, err = http.Get(server.URL + "/tasks?tag=test")
	if err != nil {
		t.Fatalf("Failed to get tasks: %v", err)
	}
	defer resp.Body.Close()

	var tasks []models.Task
	err = json.NewDecoder(resp.Body).Decode(&tasks)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
}

func TestTaskEndpoint(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a project and task
	project := models.Project{
		Name:        "Test Project",
		Description: "Test Description",
		Status:      "active",
	}

	projectData, _ := json.Marshal(project)
	resp, err := http.Post(server.URL+"/projects", "application/json", bytes.NewBuffer(projectData))
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}
	defer resp.Body.Close()

	var createdProject models.Project
	json.NewDecoder(resp.Body).Decode(&createdProject)

	task := models.Task{
		ProjectID:   createdProject.ID,
		Title:       "Test Task",
		Description: "Test Description",
		Priority:    3,
	}

	taskData, _ := json.Marshal(task)
	resp, err = http.Post(server.URL+"/tasks", "application/json", bytes.NewBuffer(taskData))
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}
	defer resp.Body.Close()

	var createdTask models.Task
	json.NewDecoder(resp.Body).Decode(&createdTask)

	// Test GET /tasks/{id}
	resp, err = http.Get(server.URL + "/tasks/" + string(rune(createdTask.ID)))
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var retrievedTask models.Task
	err = json.NewDecoder(resp.Body).Decode(&retrievedTask)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if retrievedTask.ID != createdTask.ID {
		t.Errorf("Expected task ID %d, got %d", createdTask.ID, retrievedTask.ID)
	}

	// Test PUT /tasks/{id}
	updatedTask := models.Task{
		ProjectID: createdProject.ID,
		Title:     "Updated Task",
		Priority:  5,
		Done:      true,
	}

	updatedData, _ := json.Marshal(updatedTask)
	req, _ := http.NewRequest("PUT", server.URL+"/tasks/"+string(rune(createdTask.ID)), bytes.NewBuffer(updatedData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to update task: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Test DELETE /tasks/{id}
	req, _ = http.NewRequest("DELETE", server.URL+"/tasks/"+string(rune(createdTask.ID)), nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to delete task: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestExportEndpoint(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Test JSON export
	resp, err := http.Get(server.URL + "/export?format=json")
	if err != nil {
		t.Fatalf("Failed to export JSON: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	// Test CSV export
	resp, err = http.Get(server.URL + "/export?format=csv")
	if err != nil {
		t.Fatalf("Failed to export CSV: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	contentType = resp.Header.Get("Content-Type")
	if contentType != "text/csv" {
		t.Errorf("Expected Content-Type text/csv, got %s", contentType)
	}
}

func TestStatsEndpoint(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Test stats endpoint
	resp, err := http.Get(server.URL + "/stats")
	if err != nil {
		t.Fatalf("Failed to get stats: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var stats map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&stats)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Check that stats contains expected fields
	expectedFields := []string{"total_projects", "total_tasks", "completed_tasks", "pending_tasks"}
	for _, field := range expectedFields {
		if _, exists := stats[field]; !exists {
			t.Errorf("Expected stats to contain field '%s'", field)
		}
	}
}

func TestBadgesEndpoint(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Test badges endpoint
	resp, err := http.Get(server.URL + "/badges")
	if err != nil {
		t.Fatalf("Failed to get badges: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var badges []models.Badge
	err = json.NewDecoder(resp.Body).Decode(&badges)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(badges) == 0 {
		t.Error("Expected badges list to not be empty")
	}
}
