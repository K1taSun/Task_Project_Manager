package unit

import (
	"testing"
	"time"

	"Projekt_go/internal/models"
)

func TestGenerateProjectID(t *testing.T) {
	// Reset counters before test
	models.ResetIDCounters()

	// Test ID generation
	id1 := models.GenerateProjectID()
	id2 := models.GenerateProjectID()

	if id1 != 1 {
		t.Errorf("Expected first project ID to be 1, got %d", id1)
	}
	if id2 != 2 {
		t.Errorf("Expected second project ID to be 2, got %d", id2)
	}
}

func TestGenerateTaskID(t *testing.T) {
	// Reset counters before test
	models.ResetIDCounters()

	// Test ID generation
	id1 := models.GenerateTaskID()
	id2 := models.GenerateTaskID()

	if id1 != 1 {
		t.Errorf("Expected first task ID to be 1, got %d", id1)
	}
	if id2 != 2 {
		t.Errorf("Expected second task ID to be 2, got %d", id2)
	}
}

func TestProjectExistsByID(t *testing.T) {
	// Reset counters
	models.ResetIDCounters()

	// Test non-existent project
	if models.ProjectExistsByID(1) {
		t.Error("Expected project 1 to not exist")
	}

	// Create and test existing project
	project := models.Project{
		ID:        1,
		Name:      "Test Project",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    "active",
	}
	models.SaveProject(project)

	if !models.ProjectExistsByID(1) {
		t.Error("Expected project 1 to exist")
	}
}

func TestTaskExistsByID(t *testing.T) {
	// Reset counters
	models.ResetIDCounters()

	// Test non-existent task
	if models.TaskExistsByID(1) {
		t.Error("Expected task 1 to not exist")
	}

	// Create and test existing task
	task := models.Task{
		ID:        1,
		ProjectID: 1,
		Title:     "Test Task",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Priority:  1,
	}
	models.SaveTask(task)

	if !models.TaskExistsByID(1) {
		t.Error("Expected task 1 to exist")
	}
}

func TestGetProjectsCount(t *testing.T) {
	// Reset all data
	models.ResetAllData()

	// Test empty count
	if count := models.GetProjectsCount(); count != 0 {
		t.Errorf("Expected 0 projects, got %d", count)
	}

	// Add projects and test count
	project1 := models.Project{ID: 1, Name: "Project 1", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	project2 := models.Project{ID: 2, Name: "Project 2", CreatedAt: time.Now(), UpdatedAt: time.Now()}

	models.SaveProject(project1)
	models.SaveProject(project2)

	if count := models.GetProjectsCount(); count != 2 {
		t.Errorf("Expected 2 projects, got %d", count)
	}
}

func TestGetTasksCount(t *testing.T) {
	// Reset all data
	models.ResetAllData()

	// Test empty count
	if count := models.GetTasksCount(); count != 0 {
		t.Errorf("Expected 0 tasks, got %d", count)
	}

	// Add tasks and test count
	task1 := models.Task{ID: 1, ProjectID: 1, Title: "Task 1", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	task2 := models.Task{ID: 2, ProjectID: 1, Title: "Task 2", CreatedAt: time.Now(), UpdatedAt: time.Now()}

	models.SaveTask(task1)
	models.SaveTask(task2)

	if count := models.GetTasksCount(); count != 2 {
		t.Errorf("Expected 2 tasks, got %d", count)
	}
}

func TestGetAllProjects(t *testing.T) {
	// Reset all data
	models.ResetAllData()

	// Test empty list
	projects := models.GetAllProjects()
	if len(projects) != 0 {
		t.Errorf("Expected empty projects list, got %d projects", len(projects))
	}

	// Add projects and test
	project1 := models.Project{ID: 1, Name: "Project 1", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	project2 := models.Project{ID: 2, Name: "Project 2", CreatedAt: time.Now(), UpdatedAt: time.Now()}

	models.SaveProject(project1)
	models.SaveProject(project2)

	projects = models.GetAllProjects()
	if len(projects) != 2 {
		t.Errorf("Expected 2 projects, got %d", len(projects))
	}
}

func TestGetAllTasks(t *testing.T) {
	// Reset all data
	models.ResetAllData()

	// Test empty list
	tasks := models.GetAllTasks()
	if len(tasks) != 0 {
		t.Errorf("Expected empty tasks list, got %d tasks", len(tasks))
	}

	// Add tasks and test
	task1 := models.Task{ID: 1, ProjectID: 1, Title: "Task 1", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	task2 := models.Task{ID: 2, ProjectID: 1, Title: "Task 2", CreatedAt: time.Now(), UpdatedAt: time.Now()}

	models.SaveTask(task1)
	models.SaveTask(task2)

	tasks = models.GetAllTasks()
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}
}

func TestGetProjectByID(t *testing.T) {
	// Reset all data
	models.ResetAllData()

	// Test non-existent project
	project, exists := models.GetProjectByID(1)
	if exists {
		t.Error("Expected project to not exist")
	}

	// Create and test existing project
	testProject := models.Project{
		ID:        1,
		Name:      "Test Project",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	models.SaveProject(testProject)

	project, exists = models.GetProjectByID(1)
	if !exists {
		t.Error("Expected project to exist")
	}
	if project.Name != "Test Project" {
		t.Errorf("Expected project name 'Test Project', got '%s'", project.Name)
	}
}

func TestGetTaskByID(t *testing.T) {
	// Reset all data
	models.ResetAllData()

	// Test non-existent task
	task, exists := models.GetTaskByID(1)
	if exists {
		t.Error("Expected task to not exist")
	}

	// Create and test existing task
	testTask := models.Task{
		ID:        1,
		ProjectID: 1,
		Title:     "Test Task",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	models.SaveTask(testTask)

	task, exists = models.GetTaskByID(1)
	if !exists {
		t.Error("Expected task to exist")
	}
	if task.Title != "Test Task" {
		t.Errorf("Expected task title 'Test Task', got '%s'", task.Title)
	}
}

func TestDeleteProject(t *testing.T) {
	// Reset all data
	models.ResetAllData()

	// Test deleting non-existent project
	if models.DeleteProject(1) {
		t.Error("Expected delete of non-existent project to return false")
	}

	// Create and delete project
	project := models.Project{ID: 1, Name: "Test Project", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	models.SaveProject(project)

	if !models.DeleteProject(1) {
		t.Error("Expected delete of existing project to return true")
	}

	if models.ProjectExistsByID(1) {
		t.Error("Expected project to be deleted")
	}
}

func TestDeleteTask(t *testing.T) {
	// Reset all data
	models.ResetAllData()

	// Test deleting non-existent task
	if models.DeleteTask(1) {
		t.Error("Expected delete of non-existent task to return false")
	}

	// Create and delete task
	task := models.Task{ID: 1, ProjectID: 1, Title: "Test Task", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	models.SaveTask(task)

	if !models.DeleteTask(1) {
		t.Error("Expected delete of existing task to return true")
	}

	if models.TaskExistsByID(1) {
		t.Error("Expected task to be deleted")
	}
}

func TestGetTasksByProjectID(t *testing.T) {
	// Reset all data
	models.ResetAllData()

	// Test empty project tasks
	tasks := models.GetTasksByProjectID(1)
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks for project 1, got %d", len(tasks))
	}

	// Create tasks for different projects
	task1 := models.Task{ID: 1, ProjectID: 1, Title: "Task 1", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	task2 := models.Task{ID: 2, ProjectID: 1, Title: "Task 2", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	task3 := models.Task{ID: 3, ProjectID: 2, Title: "Task 3", CreatedAt: time.Now(), UpdatedAt: time.Now()}

	models.SaveTask(task1)
	models.SaveTask(task2)
	models.SaveTask(task3)

	// Test tasks for project 1
	tasks = models.GetTasksByProjectID(1)
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks for project 1, got %d", len(tasks))
	}

	// Test tasks for project 2
	tasks = models.GetTasksByProjectID(2)
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task for project 2, got %d", len(tasks))
	}
}

func TestCreateDefaultProjectBadge(t *testing.T) {
	badge := models.CreateDefaultProjectBadge()

	if badge.Text != "New" {
		t.Errorf("Expected badge text 'New', got '%s'", badge.Text)
	}
	if badge.Type != "status" {
		t.Errorf("Expected badge type 'status', got '%s'", badge.Type)
	}
}

func TestCreatePriorityBadge(t *testing.T) {
	// Test valid priorities
	badge1 := models.CreatePriorityBadge(1)
	if badge1.Text != "Low" {
		t.Errorf("Expected priority 1 badge text 'Low', got '%s'", badge1.Text)
	}

	badge3 := models.CreatePriorityBadge(3)
	if badge3.Text != "High" {
		t.Errorf("Expected priority 3 badge text 'High', got '%s'", badge3.Text)
	}

	badge5 := models.CreatePriorityBadge(5)
	if badge5.Text != "Urgent" {
		t.Errorf("Expected priority 5 badge text 'Urgent', got '%s'", badge5.Text)
	}

	// Test invalid priority
	badgeInvalid := models.CreatePriorityBadge(10)
	if badgeInvalid.Text != "Unknown" {
		t.Errorf("Expected invalid priority badge text 'Unknown', got '%s'", badgeInvalid.Text)
	}
}

func TestUpdateProjectTimestamp(t *testing.T) {
	project := models.Project{
		ID:        1,
		Name:      "Test Project",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	oldUpdatedAt := project.UpdatedAt
	time.Sleep(1 * time.Millisecond) // Ensure time difference

	models.UpdateProjectTimestamp(&project)

	if project.UpdatedAt.Equal(oldUpdatedAt) {
		t.Error("Expected UpdatedAt to be updated")
	}
}

func TestUpdateTaskTimestamp(t *testing.T) {
	task := models.Task{
		ID:        1,
		ProjectID: 1,
		Title:     "Test Task",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	oldUpdatedAt := task.UpdatedAt
	time.Sleep(1 * time.Millisecond) // Ensure time difference

	models.UpdateTaskTimestamp(&task)

	if task.UpdatedAt.Equal(oldUpdatedAt) {
		t.Error("Expected UpdatedAt to be updated")
	}
}
